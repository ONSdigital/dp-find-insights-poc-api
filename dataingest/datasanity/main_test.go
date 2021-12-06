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
			if err := db.Raw("SELECT code FROM geo WHERE ST_Within(ST_GeomFromText('SRID=4326;POINT('|| ? || ' ' || ? ||')'),wkb_geometry::GEOMETRY) AND type_id=6", cast.ToString(tC.long), cast.ToString(tC.lat)).Scan(&results).Error; err != nil {
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
