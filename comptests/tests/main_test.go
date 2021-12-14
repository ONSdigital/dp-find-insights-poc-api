// +build comptest

package main

import (
	"context"
	"log"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/comptests"
	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/jackc/pgx/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*

This is more suggested example code for comptest with INSERT and ROLLBACK
rather than actually testing anything useful.

Also see a real test  pkg/geodata/rowquery_test.go which uses

		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")

*/

const dsn = "postgres://insights:insights@localhost:54322/censustest"

var db *gorm.DB

func init() {
	comptests.SetupDockerDB(dsn)
	model.SetupDB(dsn)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print(err)
	}

}

func TestGormExample(t *testing.T) {
	// inside transaction rolled back
	func() {
		tx := db.Begin()
		defer tx.Rollback()

		// body of test here...
		if err := tx.Save(&model.GeoType{ID: 1, Name: "EW"}).Error; err != nil {
			t.Error(err)
		}

		geot := model.GeoType{}
		tx.First(&geot)
		if geot.Name != "EW" {
			t.Fail()
		}
	}()

	// row now absent
	geot := model.GeoType{}

	db.First(&geot)

	if geot.Name == "EW" {
		t.Fail()
	}
}

func TestPgxExample(t *testing.T) {
	ctx := context.Background()
	con, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Error(err)
	}

	// inside transaction rolled back
	func() {
		tx, err := con.Begin(ctx)
		defer tx.Rollback(ctx)
		if err != nil {
			t.Error(err)
		}

		// body of test here...
		_, err = con.Exec(ctx, "INSERT INTO geo_type VALUES(1,'EW')")
		if err != nil {
			t.Error(err)
		}

		var name string

		if err = con.QueryRow(ctx, "SELECT name FROM geo_type WHERE id=$1", 1).Scan(&name); err != nil {
			t.Error(err)
		}

		if name != "EW" {
			t.Fail()
		}

	}()

	// row now absent
	var name string
	if err = con.QueryRow(ctx, "SELECT name FROM geo_type WHERE id=$1", 1).Scan(&name); err != nil {
		log.Print(err)
	}

	if name == "EW" {
		t.Fail()
	}

}
