//go:build comptest
// +build comptest

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/comptests"
	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/jackc/pgx/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = comptests.DefaultDSN

var db *gorm.DB

func init() {
	comptests.SetupDockerDB(dsn)
	model.SetupUpdateDB(dsn)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

}

type qLogger struct {
}

func (l *qLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	// uncomment me for logs
	//fmt.Printf("SQL:\n%s\nARGS:%v\n", data["sql"], data["args"])
}

func TestGetFiles(t *testing.T) {
	di := New("2011")

	di.getFiles("testdata/")

	if di.files.data[0] != "testdata/QS104EWDATA04.CSV" || di.files.meta[0] != "testdata/QS104EWMETA0.CSV" || di.files.desc[0] != "testdata/QS104EWDESC0.CSV" {
		t.Fail()
	}
}

func TestAddClassificationData(t *testing.T) {

	func() {
		var nd model.NomisDesc
		tx := db.Begin()
		defer tx.Rollback()

		di := New("2011")
		di.gdb = tx
		di.files.meta = []string{"testdata/QS104EWMETA0.CSV"}

		if foo := tx.First(&nd); !errors.Is(foo.Error, gorm.ErrRecordNotFound) {
			t.Errorf("Data wrongly present")
		}

		di.addClassificationData()

		tx.First(&nd)

		if nd.Name != "Sex" || nd.PopStat != "All usual residents" || nd.ShortNomisCode != "QS104EW" {
			t.Errorf(fmt.Sprintf("wrongly got : %#v", nd))
		}
	}()

}

func TestAddCategoryData(t *testing.T) {

	func() {
		tx := db.Begin()
		defer tx.Rollback()

		di := New("2011")
		di.gdb = tx
		di.files.meta = []string{"testdata/QS104EWMETA0.CSV"}
		di.addClassificationData()
		di.files.desc = []string{"testdata/QS104EWDESC0.CSV"}
		longToCatid := di.addCategoryData()

		if longToCatid["QS104EW0001"] == 0 || longToCatid["QS104EW0002"] == 0 {
			t.Error("data not there")
		}

		fmt.Printf("%#v\n", longToCatid)
	}()
}

func TestAddGeoGeoMetricData(t *testing.T) {
	ctx := context.Background()

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Print(err)
	}

	config.Logger = &qLogger{}
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		t.Error(err)
	}

	func() {
		tx, err := conn.Begin(ctx)
		if err != nil {
			t.Error(err)
		}
		defer tx.Rollback(ctx)

		di := New("2011")
		di.conn = conn

		conn.Exec(ctx, "INSERT INTO geo_type VALUES(4,'LAD')")
		conn.Exec(ctx, "INSERT INTO NOMIS_DESC (id,name,pop_stat,short_nomis_code,year,nomis_topic_id) VALUES (66,'Sex','All usual residents','QS104EW',2011,1)")
		conn.Exec(ctx, "INSERT INTO NOMIS_CATEGORY (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (3,66,'All categories: Sex','Count','Person','QS104EW0001',2011)")
		conn.Exec(ctx, "INSERT INTO NOMIS_CATEGORY (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (4,66,'All categories: Sex','Count','Person','QS104EW0002',2011)")

		di.files.data = []string{"testdata/QS104EWDATA04.CSV"}
		di.addGeoGeoMetricData(map[string]int32{"QS104EW0001": 3, "QS104EW0002": 4})

		var metric float64
		if err := conn.QueryRow(ctx, "SELECT metric FROM geo_metric WHERE category_id=3").Scan(&metric); err != nil {
			log.Print(err)
		}

		if metric != 92028 {
			t.Fail()
		}

	}()

}
