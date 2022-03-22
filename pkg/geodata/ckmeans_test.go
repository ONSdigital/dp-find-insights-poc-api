//go:build comptest
// +build comptest

package geodata

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/comptests"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/ONSdigital/dp-find-insights-poc-api/sentinel"
)

// ---------------------------------------------------- Helpers ----------------------------------------------------- //

func ckmeansTestSetup(t *testing.T, db *database.Database, metrics map[string]map[string][]float64) {
	// clear out any leaked data
	err := comptests.ClearDB(db)
	if err != nil {
		log.Fatal(err)
	}

	// setup data ver
	comptests.DoSQL(
		t,
		db,
		`INSERT INTO data_ver (id,created_at,updated_at,deleted_at,census_year,ver_string,source,notes,public)
		VALUES (1,'0001-01-01 00:00:00','2021-12-06 11:52:26.142808',null,2011,'2.2','Test Data','ckmeans test',true)`,
	)

	// setup topic
	comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_topic (id,top_nomis_code,name) VALUES (1,'testTopic', 'test nomis topic')"))

	// setup table
	comptests.DoSQL(t, db, "INSERT INTO nomis_desc (id,nomis_topic_id, name,pop_stat,short_nomis_code,year) VALUES (1,1,'test topic','test units','testTable',2011)")

	// setup geotypes
	geoTypeIDs := make(map[string]int)
	geoTypeID := 1
	for geoType, _ := range metrics {
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO geo_type (id,name) VALUES (%d,'%s')", geoTypeID, geoType))
		geoTypeIDs[geoType] = geoTypeID
		geoTypeID++
	}

	// setup geos
	var nGeos int
	for _, cats := range metrics {
		for _, catValues := range cats {
			nGeos = len(catValues)
			break
		}
	}
	geoIDs := make(map[string][]int)
	geoID := 1
	for geoType, geoTypeID := range geoTypeIDs {
		geoIDs[geoType] = []int{}
		for i := 1; i <= nGeos; i++ {
			comptests.DoSQL(
				t,
				db,
				fmt.Sprintf(
					`INSERT INTO geo (id,type_id,code,name,lat,long,valid,wkb_geometry,wkb_long_lat_geom)
					VALUES (%d,%d,'testGeography%d','City of Test 00%d',1,-0.1,true,null,null)`,
					geoID,
					geoTypeID,
					geoID,
					geoID,
				),
			)
			geoIDs[geoType] = append(geoIDs[geoType], geoID)
			geoID++
		}
	}

	// setup cats
	catIDs := make(map[string]map[string]int)
	catID := 1
	createdCats := make(map[string]int)
	for geoType, geoTypeMetrics := range metrics {
		catIDs[geoType] = make(map[string]int)
		for catCode, _ := range geoTypeMetrics {
			// don't setup the same category twice!
			if prevCatID, prs := createdCats[catCode]; prs {
				catIDs[geoType][catCode] = prevCatID
				continue
			}
			comptests.DoSQL(
				t,
				db,
				fmt.Sprintf(
					`INSERT INTO nomis_category (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year)
					VALUES (%d,1,'testCat%d','Count','test units','%s',2011)`,
					catID,
					catID,
					catCode,
				),
			)
			catIDs[geoType][catCode] = catID
			createdCats[catCode] = catID
			catID++
		}
	}

	// setup geo_metrics
	metricID := 1
	for geoType, geoTypeMetrics := range metrics {
		for catCode, catValues := range geoTypeMetrics {
			for i, catValue := range catValues {
				comptests.DoSQL(
					t,
					db,
					fmt.Sprintf(
						"INSERT INTO geo_metric (id,geo_id,category_id,metric,data_ver_id) VALUES (%d,%d,%d,%f,1)",
						metricID,
						geoIDs[geoType][i],
						catIDs[geoType][catCode],
						catValue,
					),
				)
				metricID++
			}
		}
	}

}

// ------------------------------------------------------ Tests ----------------------------------------------------- //

func TestCkmeansHappyPathSingleCategorySingleGeotype(t *testing.T) {
	// GIVEN the database is setup
	dsn := comptests.DefaultDSN
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")
		/*
			AND GIVEN we have seeded datapoints for one category, with values taken from
			https://github.com/simple-statistics/simple-statistics/blob/master/src/ckmeans.js#224)
		*/
		metrics := map[string]map[string][]float64{
			"LAD": {
				"category1": {-1.0, 2.0, -1.0, 2.0, 4.0, 5.0, 6.0, -1.0, 2.0, -1.0},
			},
		}
		ckmeansTestSetup(t, db, metrics)

		// WHEN we use app.CKmeans to get breakpoints for our category
		app, err := New(db, nil, 100)
		if err != nil {
			log.Fatal(err)
		}
		testK := 3
		result, err := app.CKmeans(
			context.Background(),
			2011,
			[]string{"category1"},
			[]string{"LAD"},
			testK,
			"",
		)

		// THEN we expect the breakpoints to match the example given in the original javascript repo
		wantBreaks := map[string]map[string][]float64{
			"category1": {
				"LAD": {-1.0, 2.0, 6.0},
			},
		}
		if !reflect.DeepEqual(result, wantBreaks) {
			t.Errorf("got %#v, wanted %#v", result, wantBreaks)
		}
		if err != nil {
			log.Print(err)
		}
	}()
}

func TestCkmeansHappyPathMultiCategorySingleGeotype(t *testing.T) {
	// GIVEN the database is setup
	dsn := comptests.DefaultDSN
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")
		/*
			AND GIVEN we have seeded datapoints for three category, with values adapted from
			https://github.com/simple-statistics/simple-statistics/blob/master/src/ckmeans.js#224)
		*/
		metrics := map[string]map[string][]float64{
			"LAD": {
				"category1": {-1.0, 2.0, -1.0, 2.0, 4.0, 5.0, 6.0, -1.0, 2.0, -1.0},
				"category2": {-10.0, 20.0, -10.0, 20.0, 40.0, 50.0, 60.0, -10.0, 20.0, -10.0},
				"category3": {-100.0, 200.0, -100.0, 200.0, 400.0, 500.0, 600.0, -100.0, 200.0, -100.0},
			},
		}
		ckmeansTestSetup(t, db, metrics)

		// WHEN we use app.CKmeans to get breakpoints for our category
		app, err := New(db, nil, 100)
		if err != nil {
			log.Fatal(err)
		}
		testK := 3
		result, err := app.CKmeans(
			context.Background(),
			2011,
			[]string{"category1,category2,category3"},
			[]string{"LAD"},
			testK,
			"",
		)

		// THEN we expect the breakpoints to match the example given in the original javascript repo, after adjustment
		wantBreaks := map[string]map[string][]float64{
			"category1": {
				"LAD": {-1.0, 2.0, 6.0},
			},
			"category2": {
				"LAD": {-10.0, 20.0, 60.0},
			},
			"category3": {
				"LAD": {-100.0, 200.0, 600.0},
			},
		}
		if !reflect.DeepEqual(result, wantBreaks) {
			t.Errorf("got %#v, wanted %#v", result, wantBreaks)
		}
		if err != nil {
			log.Print(err)
		}
	}()
}

func TestCkmeansHappyPathMultiCategoryMultiGeotype(t *testing.T) {
	// GIVEN the database is setup
	dsn := comptests.DefaultDSN
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")
		/*
			AND GIVEN we have seeded datapoints for three category, with values adapted from
			https://github.com/simple-statistics/simple-statistics/blob/master/src/ckmeans.js#224)
		*/
		metrics := map[string]map[string][]float64{
			"LAD": {
				"category1": {-1.0, 2.0, -1.0, 2.0, 4.0, 5.0, 6.0, -1.0, 2.0, -1.0},
				"category2": {-10.0, 20.0, -10.0, 20.0, 40.0, 50.0, 60.0, -10.0, 20.0, -10.0},
				"category3": {-100.0, 200.0, -100.0, 200.0, 400.0, 500.0, 600.0, -100.0, 200.0, -100.0},
			},
			"MSOA": {
				"category1": {-1000.0, 2000.0, -1000.0, 2000.0, 4000.0, 5000.0, 6000.0, -1000.0, 2000.0, -1000.0},
				"category2": {-10000.0, 20000.0, -10000.0, 20000.0, 40000.0, 50000.0, 60000.0, -10000.0, 20000.0, -10000.0},
				"category3": {-100000.0, 200000.0, -100000.0, 200000.0, 400000.0, 500000.0, 600000.0, -100000.0, 200000.0, -100000.0},
			},
		}
		ckmeansTestSetup(t, db, metrics)

		// WHEN we use app.CKmeans to get breakpoints for our category
		app, err := New(db, nil, 100)
		if err != nil {
			log.Fatal(err)
		}
		testK := 3
		result, err := app.CKmeans(
			context.Background(),
			2011,
			[]string{"category1,category2,category3"},
			[]string{"LAD,MSOA"},
			testK,
			"",
		)

		// THEN we expect the breakpoints to match the example given in the original javascript repo, after adjustment
		wantBreaks := map[string]map[string][]float64{
			"category1": {
				"LAD":  {-1.0, 2.0, 6.0},
				"MSOA": {-1000.0, 2000.0, 6000.0},
			},
			"category2": {
				"LAD":  {-10.0, 20.0, 60.0},
				"MSOA": {-10000.0, 20000.0, 60000.0},
			},
			"category3": {
				"LAD":  {-100.0, 200.0, 600.0},
				"MSOA": {-100000.0, 200000.0, 600000.0},
			},
		}
		if !reflect.DeepEqual(result, wantBreaks) {
			t.Errorf("got %#v, wanted %#v", result, wantBreaks)
		}
		if err != nil {
			log.Print(err)
		}
	}()
}

func TestCkmeansNoData(t *testing.T) {
	// GIVEN the database is setup
	dsn := comptests.DefaultDSN
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")

		/*
			AND GIVEN we have seeded datapoints for one category, with values taken from
			https://github.com/simple-statistics/simple-statistics/blob/master/src/ckmeans.js#224)
		*/
		metrics := map[string]map[string][]float64{"LAD": {"category1": {}}}
		ckmeansTestSetup(t, db, metrics)

		// WHEN we use app.CKmeans to get breakpoints for a category with NO DATA
		app, err := New(db, nil, 100)
		if err != nil {
			log.Fatal(err)
		}
		testK := 3
		result, err := app.CKmeans(
			context.Background(),
			2011,
			[]string{"category1"},
			[]string{"LAD"},
			testK,
			"",
		)

		// THEN we expect to receive no data
		wantBreaks := map[string]map[string][]float64{}
		if !reflect.DeepEqual(result, wantBreaks) {
			t.Errorf("got %#v, wanted %#v", result, wantBreaks)
		}
		// AND THEN we expect to receive no error
		if err != nil {
			t.Errorf("got this error = '%s', wanted nil", err)
		}
	}()
}

func TestCkmeansRatiosHappyPathSingleGeotype(t *testing.T) {
	// GIVEN the database is setup
	dsn := comptests.DefaultDSN
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")
		/*
			AND GIVEN we have seeded datapoints for our test geotype for:
				- one denominator data category
				- three numerator data categories
			Keep it simple and make the denominators all 2, and the data have three obvious order-of-magnitude breaks
		*/
		metrics := map[string]map[string][]float64{
			"LAD": {
				"denominator": {2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
				"numerator1":  {3, 598, 4, 57, 59, 60, 58, 597, 599, 600, 6, 5},
				"numerator2":  {1199, 1198, 1197, 120, 12, 118, 119, 1200, 9, 117, 10, 11},
				"numerator3":  {23, 2, 238, 240, 237, 21, 239, 2398, 2399, 2400, 24, 2397},
			},
		}
		ckmeansTestSetup(t, db, metrics)

		// WHEN we use app.CKmeansRatio to get breakpoints for all numerators / denominator
		app, err := New(db, nil, 100)
		if err != nil {
			log.Fatal(err)
		}
		testK := 3
		result, err := app.CKmeans(
			context.Background(),
			2011,
			[]string{"numerator1,numerator2,numerator3"},
			[]string{"LAD"},
			testK,
			"denominator",
		)

		// THEN we expect to get breakpoints matching the order-of-magnitude breaks in our test data
		wantBreaks := map[string]map[string][]float64{
			"numerator1": {
				"LAD": {3.0, 30.0, 300},
			},
			"numerator2": {
				"LAD": {6.0, 60.0, 600},
			},
			"numerator3": {
				"LAD": {12.0, 120.0, 1200},
			},
		}
		if !reflect.DeepEqual(result, wantBreaks) {
			t.Errorf("got %#v, wanted %#v", result, wantBreaks)
		}
		if err != nil {
			log.Print(err)
		}
	}()
}

func TestCkmeansRatiosHappyPathMultipleGeotypes(t *testing.T) {
	// GIVEN the database is setup
	dsn := comptests.DefaultDSN
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")
		/*
			AND GIVEN we have seeded datapoints for both of test geotype for:
				- one denominator data category
				- three numerator data categories
			Keep it simple and make the denominators all 2, and the data have three obvious order-of-magnitude breaks
		*/
		metrics := map[string]map[string][]float64{
			"LAD": {
				"denominator": {2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
				"numerator1":  {3, 598, 4, 57, 59, 60, 58, 597, 599, 600, 6, 5},
				"numerator2":  {1199, 1198, 1197, 120, 12, 118, 119, 1200, 9, 117, 10, 11},
				"numerator3":  {23, 2, 238, 240, 237, 21, 239, 2398, 2399, 2400, 24, 2397},
			},
			"MSOA": {
				"denominator": {2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
				"numerator1":  {480, 47, 479, 4797, 46, 478, 4800, 45, 48, 4798, 477, 4799},
				"numerator2":  {95, 958, 960, 957, 959, 94, 96, 9599, 9598, 9597, 9600, 93},
				"numerator3":  {19198, 1917, 19200, 191, 190, 1920, 189, 192, 1918, 1919, 19199, 19197},
			},
		}
		ckmeansTestSetup(t, db, metrics)

		// WHEN we use app.CKmeansRatio to get breakpoints for all numerators / denominator
		app, err := New(db, nil, 100)
		if err != nil {
			log.Fatal(err)
		}
		testK := 3
		result, err := app.CKmeans(
			context.Background(),
			2011,
			[]string{"numerator1,numerator2,numerator3"},
			[]string{"LAD,MSOA"},
			testK,
			"denominator",
		)

		// THEN we expect to get breakpoints matching the order-of-magnitude breaks in our test data
		wantBreaks := map[string]map[string][]float64{
			"numerator1": {
				"LAD":  {3.0, 30.0, 300},
				"MSOA": {24.0, 240.0, 2400},
			},
			"numerator2": {
				"LAD":  {6.0, 60.0, 600},
				"MSOA": {48, 480.0, 4800},
			},
			"numerator3": {
				"LAD":  {12.0, 120.0, 1200},
				"MSOA": {96.0, 960.0, 9600},
			},
		}
		if !reflect.DeepEqual(result, wantBreaks) {
			t.Errorf("got %#v, wanted %#v", result, wantBreaks)
		}
		if err != nil {
			log.Print(err)
		}
	}()
}

func TestCkmeansRatiosPartialDataOneCategoryMissing(t *testing.T) {
	// GIVEN the database is setup
	dsn := comptests.DefaultDSN
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")
		/*
			AND GIVEN we have seeded datapoints for our test geotype for:
				- one denominator data category
				- two numerator data categories
		*/
		metrics := map[string]map[string][]float64{
			"LAD": {
				"denominator": {2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
				"numerator1":  {3, 598, 4, 57, 59, 60, 58, 597, 599, 600, 6, 5},
				"numerator2":  {1199, 1198, 1197, 120, 12, 118, 119, 1200, 9, 117, 10, 11},
			},
		}
		ckmeansTestSetup(t, db, metrics)

		/*
			WHEN we use app.CKmeansRatio to get breakpoints for all numerators / denominator, INCLUDING ONE NUMERATOR
			THAT DOES NOT EXIST
		*/
		app, err := New(db, nil, 100)
		if err != nil {
			log.Fatal(err)
		}
		testK := 3
		result, err := app.CKmeans(
			context.Background(),
			2011,
			[]string{"numerator1,numerator2,numerator3"},
			[]string{"LAD"},
			testK,
			"denominator",
		)

		// THEN we expect to receive no data
		var wantBreaks map[string]map[string][]float64
		if !reflect.DeepEqual(result, wantBreaks) {
			t.Errorf("got %#v, wanted %#v", result, wantBreaks)
		}

		// AND THEN we expect to receive an ErrPartialContent error
		if !errors.Is(err, sentinel.ErrPartialContent) {
			t.Errorf("got this error = '%s', wanted '%s'", err, sentinel.ErrPartialContent)
		}
	}()
}

func TestCkmeansRatiosPartialDataOneCategoryPartialDataOnly(t *testing.T) {
	// GIVEN the database is setup
	dsn := comptests.DefaultDSN
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")
		/*
			AND GIVEN we have seeded datapoints for our test geotype for:
				- one denominator data category
				- full data for two numerator data categories
				- partial data for a third data category
		*/
		metrics := map[string]map[string][]float64{
			"LAD": {
				"denominator": {2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
				"numerator1":  {3, 598, 4, 57, 59, 60, 58, 597, 599, 600, 6, 5},
				"numerator2":  {1199, 1198, 1197, 120, 12, 118, 119, 1200, 9, 117, 10, 11},
				"numerator3":  {23, 2, 238, 240, 237, 21, 239, 2398, 2399, 2400},
			},
		}
		ckmeansTestSetup(t, db, metrics)

		/*
			WHEN we use app.CKmeansRatio to get breakpoints for all numerators / denominator
		*/
		app, err := New(db, nil, 100)
		if err != nil {
			log.Fatal(err)
		}
		testK := 3
		result, err := app.CKmeans(
			context.Background(),
			2011,
			[]string{"numerator1,numerator2,numerator3"},
			[]string{"LAD"},
			testK,
			"denominator",
		)

		// THEN we expect to receive no data
		var wantBreaks map[string]map[string][]float64
		if !reflect.DeepEqual(result, wantBreaks) {
			t.Errorf("got %#v, wanted %#v", result, wantBreaks)
		}

		// AND THEN we expect to receive an ErrPartialContent error
		if !errors.Is(err, sentinel.ErrPartialContent) {
			t.Errorf("got this error = '%s', wanted '%s'", err, sentinel.ErrPartialContent)
		}
	}()
}

func TestCkmeansRatiosNoData(t *testing.T) {
	// GIVEN the database is setup
	dsn := comptests.DefaultDSN
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")
		/*
			AND GIVEN we have seeded no datapoints for our test geotype
		*/
		metrics := map[string]map[string][]float64{"LAD": {}}
		ckmeansTestSetup(t, db, metrics)
		// WHEN we use app.CKmeansRatio to get breakpoints for category 1 / category 2
		app, err := New(db, nil, 100)
		if err != nil {
			log.Fatal(err)
		}
		testK := 3
		result, err := app.CKmeans(
			context.Background(),
			2011,
			[]string{"doesNotExist1,doesNotExist2"},
			[]string{"LAD"},
			testK,
			"doesNotExist3",
		)

		// THEN we expect to receive no data
		wantBreaks := map[string]map[string][]float64{}
		if !reflect.DeepEqual(result, wantBreaks) {
			t.Errorf("got %#v, wanted %#v", result, wantBreaks)
		}

		// AND THEN we expect to receive no error
		if err != nil {
			t.Errorf("got this error = '%s', wanted nil", err)
		}
	}()
}

func TestCkmeansArgParsingAndValidation(t *testing.T) {
	// GIVEN the database is setup
	dsn := comptests.DefaultDSN
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")
		/*
			AND GIVEN we have seeded datapoints for both of test geotype for:
				- one denominator data category
				- three numerator data categories
			Keep it simple and make the denominators all 2, and the data have three obvious order-of-magnitude breaks
		*/
		metrics := map[string]map[string][]float64{
			"LAD": {
				"denominator": {2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
				"numerator1":  {3, 598, 4, 57, 59, 60, 58, 597, 599, 600, 6, 5},
				"numerator2":  {1199, 1198, 1197, 120, 12, 118, 119, 1200, 9, 117, 10, 11},
				"numerator3":  {23, 2, 238, 240, 237, 21, 239, 2398, 2399, 2400, 24, 2397},
			},
			"MSOA": {
				"denominator": {2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
				"numerator1":  {480, 47, 479, 4797, 46, 478, 4800, 45, 48, 4798, 477, 4799},
				"numerator2":  {95, 958, 960, 957, 959, 94, 96, 9599, 9598, 9597, 9600, 93},
				"numerator3":  {19198, 1917, 19200, 191, 190, 1920, 189, 192, 1918, 1919, 19199, 19197},
			},
		}
		ckmeansTestSetup(t, db, metrics)

		app, err := New(db, nil, 100)
		if err != nil {
			log.Fatal(err)
		}
		testK := 3

		// WHEN we use app.CKmeansRatio to get breakpoints for all numerators / denominator, using various versions of
		// calling it...
		for _, argset := range []map[string][]string{
			{
				"cat":     []string{"numerator1,numerator2,numerator3"},
				"geotype": []string{"LAD,MSOA"},
			},
			{
				"cat":     []string{"numerator1", "numerator2", "numerator3"},
				"geotype": []string{"LAD", "MSOA"},
			},
			{
				"cat":     []string{"numerator1,numerator2,numerator3"},
				"geotype": []string{"lad,msoa"},
			},
		} {
			result, err := app.CKmeans(
				context.Background(),
				2011,
				argset["cat"],
				argset["geotype"],
				testK,
				"denominator",
			)

			// THEN we expect to get breakpoints matching the order-of-magnitude breaks in our test data, in all cases
			wantBreaks := map[string]map[string][]float64{
				"numerator1": {
					"LAD":  {3.0, 30.0, 300},
					"MSOA": {24.0, 240.0, 2400},
				},
				"numerator2": {
					"LAD":  {6.0, 60.0, 600},
					"MSOA": {48, 480.0, 4800},
				},
				"numerator3": {
					"LAD":  {12.0, 120.0, 1200},
					"MSOA": {96.0, 960.0, 9600},
				},
			}
			if !reflect.DeepEqual(result, wantBreaks) {
				t.Errorf("got %#v, wanted %#v", result, wantBreaks)
			}
			if err != nil {
				log.Print(err)
			}
		}

		// AND WHEN we try to call using ranges, we get an error
		result, err := app.CKmeans(
			context.Background(),
			2011,
			[]string{"numerator1...numerator3"},
			[]string{"LAD,MSOA"},
			testK,
			"denominator",
		)

		// THEN we expect to get breakpoints matching the order-of-magnitude breaks in our test data, in all cases
		var wantBreaks map[string]map[string][]float64
		if !reflect.DeepEqual(result, wantBreaks) {
			t.Errorf("got %#v, wanted %#v", result, wantBreaks)
		}

		// AND THEN we expect to receive an ErrInvalidParams error
		if !errors.Is(err, sentinel.ErrInvalidParams) {
			t.Errorf("got this error = '%s', wanted '%s'", err, sentinel.ErrInvalidParams)
		}
	}()
}
