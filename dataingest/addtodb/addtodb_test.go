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
	model.SetupDBOnceOnly(dsn)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

}

type qLogger struct {
}

func (l *qLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	//spew.Dump(data)
	//if level == pgx.LogLevelInfo && msg == "Query" {
	fmt.Printf("SQL:\n%s\nARGS:%v\n", data["sql"], data["args"])
	//}
}

func TestGetFiles(t *testing.T) {
	di := New("2011")

	di.getFiles("testdata/")

	if di.files.data[0] != "testdata/QS104EWDATA04.CSV" || di.files.meta[0] != "testdata/QS104EWMETA0.CSV" || di.files.desc[0] != "testdata/QS104EWDESC0.CSV" {
		t.Fail()
	}
}

func TestAddMetaTables(t *testing.T) {

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

		di.addMetaTables()

		tx.First(&nd)

		if nd.Name != "Sex" || nd.PopStat != "All usual residents" || nd.ShortNomisCode != "QS104EW" {
			t.Errorf(fmt.Sprintf("wrongly got : %#v", nd))
		}
	}()

}

func TestAddDiscTables(t *testing.T) {

	func() {
		tx := db.Begin()
		defer tx.Rollback()

		di := New("2011")
		di.gdb = tx
		di.files.meta = []string{"testdata/QS104EWMETA0.CSV"}
		di.addMetaTables()
		di.files.desc = []string{"testdata/QS104EWDESC0.CSV"}
		longToCatid := di.addDiscTables()

		if longToCatid["QS104EW0001"] == 0 || longToCatid["QS104EW0002"] == 0 {
			t.Error("data not there")
		}

		fmt.Printf("%#v\n", longToCatid)
	}()
}

/*
func TestAddDataTables(t *testing.T) {
	ctx := context.Background()

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Print(err)
	}

	config.Logger = &qLogger{}

	con, err := pgx.ConnectConfig(ctx, config)

	//	con, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Error(err)
	}

	func() {
		tx, err := con.Begin(ctx)
		defer tx.Rollback(ctx)
		if err != nil {
			t.Error(err)
		}

		di := New("2011")
		di.conn = con
		di.gdb = db // none rolling
		di.files.data = []string{"testdata/QS104EWDATA04.CSV"}

		di.createGeoTypes() // need to check this rolls back XXX
		var longToCatid map[string]int32
		di.addDataTables(longToCatid)

		// geo & geo_metric

		/*
			if foo := tx.First(&nd); !errors.Is(foo.Error, gorm.ErrRecordNotFound) {
				t.Errorf("Data wrongly present")
			}

			di.addMetaTables()

			tx.First(&nd)

			if nd.Name != "Sex" || nd.PopStat != "All usual residents" || nd.ShortNomisCode != "QS104EW" {
				t.Errorf(fmt.Sprintf("wrongly got : %#v", nd))
			}
	}()

}
*/
