//go:build datasanity
// +build datasanity

package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/spf13/cast"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(postgres.Open(database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

}

func TestSpacialLatLog(t *testing.T) {
	testCases := []struct {
		lat          float64
		long         float64
		desc         string
		expectedCode string
	}{
		{
			lat:          51.895167,
			long:         1.4805,
			desc:         "Sealand",
			expectedCode: "",
		},
		{
			lat:          53.06856,
			long:         -4.076072,
			desc:         "Snowdon Peak",
			expectedCode: "W01000118",
		},
		{
			lat:          51.476852,
			long:         -0.000500,
			desc:         "The Royal Observatory Greenwich",
			expectedCode: "E01001642",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			results := []model.Geo{}
			if err := db.Raw(`
			SELECT code FROM geo 
			WHERE ST_Within(ST_GeomFromText('POINT('|| ? || ' ' || ? ||')',4326),wkb_geometry::GEOMETRY) 
			AND type_id=6
			`, cast.ToString(tC.long), cast.ToString(tC.lat)).Scan(&results).Error; err != nil {
				t.Error(err)
			}

			fmt.Printf("%#v\n", results)

			if len(results) > 1 {
				t.Errorf("more than one row")
			}

			if len(results) == 1 {
				if results[0].Code != tC.expectedCode {
					t.Errorf("got %s code", results[0].Code)
				}
			} else if tC.expectedCode != "" {
				t.Errorf("got unexpected result")
			}

		})
	}
}

func TestWelshAbsent(t *testing.T) {
	results := []model.NomisCategory{}
	if err := db.Raw(`
	SELECT FROM nomis_category 
	WHERE category_name 
	LIKE '%Cyfradd%ddeiliadaeth%'
	`).Scan(&results).Error; err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", results)

	if len(results) != 0 {
		t.Errorf("got unexpected row(s)")
	}
}

func TestMsoaDataPresent(t *testing.T) {
	var count int
	if err := db.Raw(`
	SELECT count(*) 
	FROM geo_metric, geo
	WHERE geo_metric.geo_id=geo.id 
	AND geo.type_id=5
	`).Scan(&count).Error; err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", count)

	if count == 0 {
		t.Error("MSOA data not there")
	}
}

func TestMsoaCodesPresent(t *testing.T) {
	var count int
	if err := db.Raw(`
	SELECT count(*) 
	FROM geo
	WHERE type_id=5
	`).Scan(&count).Error; err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", count)

	if count == 0 {
		t.Error("MSOA codes not there")
	}
}

func TestAllGeoNamed(t *testing.T ) {
	var count int
	if err := db.Raw(`
	SELECT count(*) 
	FROM geo
	WHERE name='NA' AND valid=true
	`).Scan(&count).Error; err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", count)

	if count != 0 {
		t.Error("not all geo named")
	}
}

func TestGeomUKBbox(t *testing.T) {
	var codes []string
	// UK like bbox
	if err := db.Raw(`
	SELECT code FROM geo 
	WHERE NOT geo.wkb_geometry && ST_GeomFromText( 'MULTIPOINT( -7.57 49.92, 1.76 58.64)', 4326)
	`).Scan(&codes).Error; err != nil {
		t.Error(err)
	}

	if len(codes) > 0 {
		t.Errorf("got unexpected row(s) %v", codes)
	}
}

func TestLatLongGeom(t *testing.T) {
	var codes []string
	// UK like bbox
	if err := db.Raw(`
	SELECT code FROM geo 
	WHERE NOT geo.wkb_long_lat_geom && ST_GeomFromText( 'MULTIPOINT( -7.57 49.92, 1.76 58.64)', 4326)`).Scan(&codes).Error; err != nil {
		t.Error(err)
	}

	if len(codes) > 0 {
		t.Errorf("got unexpected row(s) %v", codes)
	}
}

// some value queries - will break with different data than 2011
func TestSomeValues(t *testing.T) {
	metric := model.GeoMetric{}
	db.First(&metric)

	if metric.Metric != 45496780.0 {
		t.Errorf("got %f", metric.Metric)
	}

	geo := model.Geo{}
	db.First(&geo)

	if geo.Name != "England and Wales" {
		t.Errorf("got %s", geo.Name)
	}

}

// bulk data long nomis codes have different length to API ones!
func TestLongNomisCode(t *testing.T) {
	var length []int
	// UK like bbox
	if err := db.Raw(`
    SELECT DISTINCT(LENGTH(long_nomis_code)) 
	FROM nomis_category
	`).Scan(&length).Error; err != nil {
		t.Error(err)
	}

	if len(length) > 1 {
		t.Errorf("got unexpected row(s) %v", length)
	}
}

// check short nomis
func TestShortNomisCode(t *testing.T) {

	var got []string
	if err := db.Raw(`
    SELECT short_nomis_code 
	FROM nomis_desc 
	ORDER BY short_nomis_code ASC
	`).Scan(&got).Error; err != nil {
		t.Error(err)
	}

	expected := expectedCodes()

	// these exist in addtodb/2i.txt but not a recent snapshot of v4 sheet
	for _, e := range got {
		if !elemInSlice(e, expected) {
			fmt.Printf("%s extra\n", e)
		}
	}

	for _, e := range expected {
		if !elemInSlice(e, got) {
			t.Errorf("%s missing\n", e)
		}
	}

}

func elemInSlice(e string, ss []string) (in bool) {
	for _, s := range ss {
		if e == s {
			return true
		}
	}

	return false
}

// from Viv's v4 google spreadsheet at time of modification
func expectedCodes() []string {
	return []string{
		"KS103EW",
		"KS202EW",
		"KS206EW",
		"KS207WA",
		"KS608EW",
		"QS101EW",
		"QS103EW",
		"QS104EW",
		"QS113EW",
		"QS119EW",
		"QS201EW",
		"QS202EW",
		"QS203EW",
		"QS208EW",
		"QS301EW",
		"QS302EW",
		"QS303EW",
		"QS402EW",
		"QS403EW",
		"QS406EW",
		"QS411EW",
		"QS415EW",
		"QS416EW",
		"QS501EW",
		"QS601EW",
		"QS604EW",
		"QS605EW",
		"QS701EW",
		"QS702EW",
		"QS803EW"}

}
