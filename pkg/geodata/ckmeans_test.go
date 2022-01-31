//go:build comptest
// +build comptest

package geodata

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/comptests"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
)

func TestCkmeansHappyPath(t *testing.T) {
	// GIVEN the database is setup
	dsn := comptests.DefaultDSN
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")

		// clear out pre-pop data
		err := comptests.ClearDB(db)
		if err != nil {
			log.Fatal(err)
		}

		// AND GIVEN we have seeded datapoints for a single data category
		testGeotype := "TestGeoType"
		testTopic := "TestTopic"
		testTable := testTopic + "TestTable"
		testCat := testTable + "1"
		testK := 3
		comptests.DoSQL(t, db, "INSERT INTO data_ver (id,created_at,updated_at,deleted_at,census_year,ver_string,source,notes,public) VALUES (1,'0001-01-01 00:00:00','2021-12-06 11:52:26.142808',null,2011,'2.2','Test Data','ckmeans test',true)")
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO geo_type (id,name) VALUES (1,'%s')", testGeotype))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_topic (id,top_nomis_code,name) VALUES (1,'%s', 'test nomis topic')", testTopic))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_desc (id,nomis_topic_id, name,pop_stat,short_nomis_code,year) VALUES (1,1,'test topic','test units','%s',2011)", testTable))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_category (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (1,1,'test metric','Count','test units','%s',2011)", testCat))

		// seed data for ckmeans test (values taken from https://github.com/simple-statistics/simple-statistics/blob/master/src/ckmeans.js#224)
		geoMetricValues := []float64{-1.0, 2.0, -1.0, 2.0, 4.0, 5.0, 6.0, -1.0, 2.0, -1.0}
		for i, geoMetricValue := range geoMetricValues {
			id := i + 1
			geographySQL := fmt.Sprintf(
				"INSERT INTO geo (id,type_id,code,name,lat,long,valid,wkb_geometry,wkb_long_lat_geom) VALUES (%d,1,'testGeograpy%d','City of Test 00%d',1,-0.1,true,null,null)",
				id,
				id,
				id,
			)
			comptests.DoSQL(t, db, geographySQL)
			metricSQL := fmt.Sprintf(
				"INSERT INTO geo_metric (id,geo_id,category_id,metric,data_ver_id) VALUES (%d,%d,1,%f,1)",
				id,
				id,
				geoMetricValue,
			)
			comptests.DoSQL(t, db, metricSQL)
		}

		// WHEN we use app.CKmeans to get breakpoints for our category
		app, err := New(db, 100)
		if err != nil {
			log.Fatal(err)
		}

		result, err := app.CKmeans(
			context.Background(),
			2011,
			testCat,
			testGeotype,
			testK,
		)

		// THEN we expect the breakpoints to match the example given at https://github.com/simple-statistics/simple-statistics/blob/master/src/ckmeans.js#224
		want := []float64{
			-1.0,
			2.0,
			6.0,
		}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("got %#v", result)
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

		// clear out pre-pop data
		err := comptests.ClearDB(db)
		if err != nil {
			log.Fatal(err)
		}

		// AND GIVEN we NOT have seeded datapoints for any categories categories
		testGeotype := "TestGeoType"
		testTopic := "TestTopic"
		testTable := testTopic + "TestTable"
		testCat := testTable + "1"
		testK := 3

		// WHEN we use app.CKmeans to get breakpoints for a category with NO DATA
		app, err := New(db, 100)
		if err != nil {
			log.Fatal(err)
		}

		result, err := app.CKmeans(
			context.Background(),
			2011,
			testCat,
			testGeotype,
			testK,
		)

		// THEN we expect to receive no data
		var wantData []float64
		if !reflect.DeepEqual(result, wantData) {
			t.Errorf("got %#v", result)
		}

		// AND THEN we expect to receive an ErrNoContent error
		if !reflect.DeepEqual(err, ErrNoContent) {
			t.Errorf("got this error = '%s', wanted '%s'", err, ErrNoContent)
		}
	}()
}

func TestCkmeansRatiosHappyPath(t *testing.T) {
	// GIVEN the database is setup
	dsn := comptests.DefaultDSN
	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")

		// clear out pre-pop data
		err := comptests.ClearDB(db)
		if err != nil {
			log.Fatal(err)
		}

		// AND GIVEN we have seeded datapoints for two data categories
		testGeotype := "TestGeoType"
		testTopic := "TestTopic"
		testTable := testTopic + "TestTable"
		testCat1 := testTable + "1"
		testCat2 := testTable + "2"
		testK := 3
		comptests.DoSQL(t, db, "INSERT INTO data_ver (id,created_at,updated_at,deleted_at,census_year,ver_string,source,notes,public) VALUES (1,'0001-01-01 00:00:00','2021-12-06 11:52:26.142808',null,2011,'2.2','Test Data','ckmeans test',true)")
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO geo_type (id,name) VALUES (1,'%s')", testGeotype))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_topic (id,top_nomis_code,name) VALUES (1,'%s', 'test nomis topic')", testTopic))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_desc (id,nomis_topic_id, name,pop_stat,short_nomis_code,year) VALUES (1,1,'test topic','test units','%s',2011)", testTable))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_category (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (1,1,'test metric 1','Count','test units','%s',2011)", testCat1))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_category (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (2,1,'test metric 2','Count','test units','%s',2011)", testCat2))

		// seed data for ckmeans test (metric1/metric2 gives same distribution as example from https://github.com/simple-statistics/simple-statistics/blob/master/src/ckmeans.js#224)
		geoMetric1Values := []float64{8, 35, 8, 35, 63, 80, 99, 8, 35, 8}
		geoMetric2Values := []float64{2, 5, 2, 5, 7, 8, 9, 2, 5, 2}
		for i, geoMetricValue := range geoMetric1Values {
			id := i + 1
			geographySQL := fmt.Sprintf(
				"INSERT INTO geo (id,type_id,code,name,lat,long,valid,wkb_geometry,wkb_long_lat_geom) VALUES (%d,1,'testGeograpy%d','City of Test 00%d',1,-0.1,true,null,null)",
				id,
				id,
				id,
			)
			comptests.DoSQL(t, db, geographySQL)
			metricSQL := fmt.Sprintf(
				"INSERT INTO geo_metric (id,geo_id,category_id,metric,data_ver_id) VALUES (%d,%d,1,%f,1)",
				id,
				id,
				geoMetricValue,
			)
			comptests.DoSQL(t, db, metricSQL)
		}
		for i, geoMetricValue := range geoMetric2Values {
			geoID := i + 1
			id := i + len(geoMetric2Values) + 1
			metricSQL := fmt.Sprintf(
				"INSERT INTO geo_metric (id,geo_id,category_id,metric,data_ver_id) VALUES (%d,%d,2,%f,1)",
				id,
				geoID,
				geoMetricValue,
			)
			comptests.DoSQL(t, db, metricSQL)
		}

		// WHEN we use app.CKmeansRatio to get breakpoints for category 1 / category 2
		app, err := New(db, 100)
		if err != nil {
			log.Fatal(err)
		}

		result, err := app.CKmeansRatio(
			context.Background(),
			2011,
			testCat1,
			testCat2,
			testGeotype,
			testK,
		)

		// THEN we expect the breakpoints to match the example given at https://github.com/simple-statistics/simple-statistics/blob/master/src/ckmeans.js#224
		// (with 5 added to each point to avoid issues producing negative ratios)
		want := []float64{
			4.0,
			7.0,
			11.0,
		}

		if !reflect.DeepEqual(result, want) {
			t.Errorf("got %#v", result)
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

		// clear out pre-pop data
		err := comptests.ClearDB(db)
		if err != nil {
			log.Fatal(err)
		}

		// AND GIVEN we have seeded datapoints for ONE data category
		testGeotype := "TestGeoType"
		testTopic := "TestTopic"
		testTable := testTopic + "TestTable"
		testCat1 := testTable + "1"
		testCat2 := testTable + "2"
		testK := 3
		comptests.DoSQL(t, db, "INSERT INTO data_ver (id,created_at,updated_at,deleted_at,census_year,ver_string,source,notes,public) VALUES (1,'0001-01-01 00:00:00','2021-12-06 11:52:26.142808',null,2011,'2.2','Test Data','ckmeans test',true)")
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO geo_type (id,name) VALUES (1,'%s')", testGeotype))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_topic (id,top_nomis_code,name) VALUES (1,'%s', 'test nomis topic')", testTopic))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_desc (id,nomis_topic_id, name,pop_stat,short_nomis_code,year) VALUES (1,1,'test topic','test units','%s',2011)", testTable))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_category (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (1,1,'test metric 1','Count','test units','%s',2011)", testCat1))

		// seed data for ckmeans test (metric1/metric2 gives same distribution as example from https://github.com/simple-statistics/simple-statistics/blob/master/src/ckmeans.js#224)
		geoMetric1Values := []float64{8, 35, 8, 35, 63, 80, 99, 8, 35, 8}
		for i, geoMetricValue := range geoMetric1Values {
			id := i + 1
			geographySQL := fmt.Sprintf(
				"INSERT INTO geo (id,type_id,code,name,lat,long,valid,wkb_geometry,wkb_long_lat_geom) VALUES (%d,1,'testGeograpy%d','City of Test 00%d',1,-0.1,true,null,null)",
				id,
				id,
				id,
			)
			comptests.DoSQL(t, db, geographySQL)
			metricSQL := fmt.Sprintf(
				"INSERT INTO geo_metric (id,geo_id,category_id,metric,data_ver_id) VALUES (%d,%d,1,%f,1)",
				id,
				id,
				geoMetricValue,
			)
			comptests.DoSQL(t, db, metricSQL)
		}

		// WHEN we use app.CKmeansRatio to get breakpoints for category 1 / category 2
		app, err := New(db, 100)
		if err != nil {
			log.Fatal(err)
		}

		result, err := app.CKmeansRatio(
			context.Background(),
			2011,
			testCat1,
			testCat2,
			testGeotype,
			testK,
		)

		// THEN we expect to receive no data
		var wantData []float64
		if !reflect.DeepEqual(result, wantData) {
			t.Errorf("got %#v", result)
		}

		// AND THEN we expect to receive an ErrPartialContent error
		if !reflect.DeepEqual(err, ErrPartialContent) {
			t.Errorf("got this error = '%s', wanted '%s'", err, ErrPartialContent)
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

		// clear out pre-pop data
		err := comptests.ClearDB(db)
		if err != nil {
			log.Fatal(err)
		}

		// AND GIVEN we have seeded complete datapoints for one data categories, but only partial for a second
		testGeotype := "TestGeoType"
		testTopic := "TestTopic"
		testTable := testTopic + "TestTable"
		testCat1 := testTable + "1"
		testCat2 := testTable + "2"
		testK := 3
		comptests.DoSQL(t, db, "INSERT INTO data_ver (id,created_at,updated_at,deleted_at,census_year,ver_string,source,notes,public) VALUES (1,'0001-01-01 00:00:00','2021-12-06 11:52:26.142808',null,2011,'2.2','Test Data','ckmeans test',true)")
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO geo_type (id,name) VALUES (1,'%s')", testGeotype))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_topic (id,top_nomis_code,name) VALUES (1,'%s', 'test nomis topic')", testTopic))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_desc (id,nomis_topic_id, name,pop_stat,short_nomis_code,year) VALUES (1,1,'test topic','test units','%s',2011)", testTable))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_category (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (1,1,'test metric 1','Count','test units','%s',2011)", testCat1))
		comptests.DoSQL(t, db, fmt.Sprintf("INSERT INTO nomis_category (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (2,1,'test metric 2','Count','test units','%s',2011)", testCat2))

		// seed data for ckmeans test (metric1/metric2 gives same distribution as example from https://github.com/simple-statistics/simple-statistics/blob/master/src/ckmeans.js#224)
		geoMetric1Values := []float64{8, 35, 8, 35, 63, 80, 99, 8, 35, 8}
		geoMetric2Values := []float64{2, 5, 2, 5, 7, 8, 9, 2, 5, 2}
		for i, geoMetricValue := range geoMetric1Values {
			id := i + 1
			geographySQL := fmt.Sprintf(
				"INSERT INTO geo (id,type_id,code,name,lat,long,valid,wkb_geometry,wkb_long_lat_geom) VALUES (%d,1,'testGeograpy%d','City of Test 00%d',1,-0.1,true,null,null)",
				id,
				id,
				id,
			)
			comptests.DoSQL(t, db, geographySQL)
			metricSQL := fmt.Sprintf(
				"INSERT INTO geo_metric (id,geo_id,category_id,metric,data_ver_id) VALUES (%d,%d,1,%f,1)",
				id,
				id,
				geoMetricValue,
			)
			comptests.DoSQL(t, db, metricSQL)
		}
		// only write half the values for the second metric
		for i, geoMetricValue := range geoMetric2Values {
			if i > len(geoMetric2Values)/2 {
				break
			}
			geoID := i + 1
			id := i + len(geoMetric2Values) + 1
			metricSQL := fmt.Sprintf(
				"INSERT INTO geo_metric (id,geo_id,category_id,metric,data_ver_id) VALUES (%d,%d,2,%f,1)",
				id,
				geoID,
				geoMetricValue,
			)
			comptests.DoSQL(t, db, metricSQL)
		}

		// WHEN we use app.CKmeansRatio to get breakpoints for category 1 / category 2
		app, err := New(db, 100)
		if err != nil {
			log.Fatal(err)
		}

		result, err := app.CKmeansRatio(
			context.Background(),
			2011,
			testCat1,
			testCat2,
			testGeotype,
			testK,
		)

		// THEN we expect to receive no data
		var wantData []float64
		if !reflect.DeepEqual(result, wantData) {
			t.Errorf("got %#v", result)
		}

		// AND THEN we expect to receive an ErrPartialContent error
		if !reflect.DeepEqual(err, ErrPartialContent) {
			t.Errorf("got this error = '%s', wanted '%s'", err, ErrPartialContent)
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

		// clear out pre-pop data
		err := comptests.ClearDB(db)
		if err != nil {
			log.Fatal(err)
		}

		// AND GIVEN we NOT have seeded datapoints for any categories categories
		testGeotype := "TestGeoType"
		testTopic := "TestTopic"
		testTable := testTopic + "TestTable"
		testCat1 := testTable + "1"
		testCat2 := testTable + "2"
		testK := 3

		// WHEN we use app.CKmeansRatio to get breakpoints for category 1 / category 2
		app, err := New(db, 100)
		if err != nil {
			log.Fatal(err)
		}

		result, err := app.CKmeansRatio(
			context.Background(),
			2011,
			testCat1,
			testCat2,
			testGeotype,
			testK,
		)

		// THEN we expect to receive no data
		var wantData []float64
		if !reflect.DeepEqual(result, wantData) {
			t.Errorf("got %#v", result)
		}

		// AND THEN we expect to receive an ErrNoContent error
		if !reflect.DeepEqual(err, ErrNoContent) {
			t.Errorf("got this error = '%s', wanted '%s'", err, ErrNoContent)
		}
	}()
}
