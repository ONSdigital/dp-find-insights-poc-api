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

func TestMsoaDataAbsent(t *testing.T) {
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

	if count != 0 {
		t.Error("got unexpected row(s)")
	}
}

func TestMsoaCodesAbsent(t *testing.T) {
	var count int
	if err := db.Raw(`
	SELECT count(*) 
	FROM geo
	WHERE type_id=5
	`).Scan(&count).Error; err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", count)

	if count != 0 {
		t.Error("got unexpected row(s)")
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
