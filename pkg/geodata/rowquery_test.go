// +build comptest

package geodata

import (
	"context"
	"log"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/comptests"
	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
)

// TODO Check empty
func TestRowQuery(t *testing.T) {
	const dsn = comptests.DefaultDSN
	comptests.SetupDockerDB(dsn)
	model.SetupDBOnceOnly(dsn)

	db, err := database.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		db.DB().Exec("BEGIN")
		defer db.DB().Exec("ROLLBACK")

		// setup database
		comptests.DoSQL(t, db, "INSERT INTO GEO_TYPE (id,name) VALUES (6,'LSOA')")
		comptests.DoSQL(t, db, "INSERT INTO GEO (id,type_id,code,name,lat,long,valid,wkb_geometry,wkb_long_lat_geom) VALUES (7562,6,'E01000001','City of London 001A',51.5181,-0.09706,true,null,null)")
		comptests.DoSQL(t, db, "INSERT INTO DATA_VER (id,created_at,updated_at,deleted_at,census_year,ver_string,source,notes,public) VALUES (2,'0001-01-01 00:00:00','2021-12-06 11:52:26.142808',null,2011,'2.2','Nomis Bulk API','Release date 12/02/2013 Revised 17/01/2014',true)")
		comptests.DoSQL(t, db, "INSERT INTO nomis_desc (id,nomis_topic_id,name,pop_stat,short_nomis_code,year) VALUES (63,1,'Population density','All usual residents; Area (Hectares)','QS102EW',2011)")
		comptests.DoSQL(t, db, "INSERT INTO NOMIS_CATEGORY (id,nomis_desc_id,category_name,measurement_unit,stat_unit,long_nomis_code,year) VALUES (1293,63,'All usual residents','Count','Person','QS102EW0001',2011)")
		comptests.DoSQL(t, db, "INSERT INTO geo_metric (id,geo_id,category_id,metric,data_ver_id) VALUES (54575893,7562,1293,1465.0,1)")

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
			[]string{"geography_code"},
		)

		if result != "geography_code,QS102EW0001\nE01000001,1465\n" {
			t.Errorf("got %#v", result)
		}

		if err != nil {
			log.Print(err)
		}
	}()

}
