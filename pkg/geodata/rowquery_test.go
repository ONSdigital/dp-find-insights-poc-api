// +build comptest

package geodata

import (
	"context"
	"flag"
	"log"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/comptests"
	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
)

const dsn = "postgres://insights:insights@localhost:54323/censustest"

// passing -args -kill=false to the test will kill docker postgres
var kill = flag.Bool("kill", false, "docker kill postgres")

// TODO Check empty
func TestRowQuery(t *testing.T) {
	comptests.SetupDockerDB(dsn)
	model.SetupDB(dsn)

	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")

		// setup database
		doSQL(t, db, "INSERT INTO GEO_TYPE (id,name) VALUES (6,'LSOA')")
		doSQL(t, db, "INSERT INTO GEO (id,type_id,code,name,lat,long,valid,wkb_geometry,wkb_long_lat_geom) VALUES (7562,6,'E01000001','City of London 001A',51.5181,-0.09706,true,'0103000020E6100000010000000600000099B7D26A2C41B8BFCABAA1E4A2C249403E6F5814C06FB8BF4B7A3CFEF9C149407A98D8707087B9BF11B32F7F25C249403FBA9D552F37B9BF7CB1250DA1C24940ECB58D3C65E6B8BF8C2B844AC3C2494099B7D26A2C41B8BFCABAA1E4A2C24940','0101000020E610000045F0BF95ECD8B8BF5F07CE1951C24940')")
		doSQL(t, db, "INSERT INTO DATA_VER (id,created_at,updated_at,deleted_at,census_year,ver_string,source,notes,public) VALUES (1,'0001-01-01 00:00:00','2021-12-06 11:52:26.142808',null,2011,'2.2','Nomis Bulk API','Release date 12/02/2013 Revised 17/01/2014',true)")
		doSQL(t, db, "INSERT INTO nomis_desc (id,name,pop_stat,short_nomis_code,year) VALUES (63,'Population density','All usual residents; Area (Hectares)','QS102EW',2011)")
		doSQL(t, db, "INSERT INTO NOMIS_CATEGORY (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (1293,63,'All usual residents','Count','Person','QS102EW0001',2011)")
		doSQL(t, db, "INSERT INTO geo_metric (id,geo_id,category_id,metric,data_ver_id) VALUES (54575893,7562,1293,1465.0,1)")

		// body of test here...
		app, err := New(db, 100)
		if err != nil {
			log.Fatal(err)
		}

		result, err := app.rowQuery(
			context.Background(),
			[]string{"E01000001"},   // geos
			[]string{"LSOA"},        // geotypes
			[]string{"QS102EW0001"}, // cats
		)

		if result != "geography_code,QS102EW0001\nE01000001,1465\n" {
			t.Errorf("got %#v", result)
		}

		if err != nil {
			log.Print(err)
		}
	}()

}

func TestDockerKill(t *testing.T) {
	// Dummy test to optionally take down docker
	if *kill {
		comptests.KillDockerDB()
	}
}

func doSQL(t *testing.T, db *database.Database, sql string) {
	_, err := db.DB().Exec(sql)
	if err != nil {
		t.Fatal(err)
	}
}
