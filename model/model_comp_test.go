package geobb

import (
	"log"
	"reflect"
	"testing"

	"github.com/ONSdigital/dp-find-insights-poc-api/comptests"
	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// comptest

// This should probably really be under "model" namespace but there is an
// import loop!

const dsn = comptests.DefaultDSN

var gdb *gorm.DB

func init() {
	comptests.SetupDockerDB(dsn)
	model.SetupUpdateDB(dsn)

	var err error
	gdb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Print(err)
	}
}

func TestGeometryFetchSave(t *testing.T) {
	// inside transaction rolled back
	func() {
		tx := gdb.Begin()
		defer tx.Rollback()

		if err := tx.Exec("INSERT INTO geo_type (id,name) VALUES (6,'LSOA')").Error; err != nil {
			t.Fatalf(err.Error())
		}

		if err := tx.Exec("INSERT INTO geo (id,type_id,code,name,lat,long,valid,wkb_geometry,wkb_long_lat_geom) VALUES (7563,6,'E01000002','City of London 001B',51.51868,-0.09197,true,'0103000020E6100000010000000500000091C765C7EA8DB6BFDAFE241B7CC249403E6F5814C06FB8BF4B7A3CFEF9C1494099B7D26A2C41B8BFCABAA1E4A2C24940A65B75599DBDB7BFA6358AE0BCC2494091C765C7EA8DB6BFDAFE241B7CC24940','0101000020E6100000088F368E588BB7BF6E8B321B64C24940')").Error; err != nil {
			t.Errorf(err.Error())
		}

		var g model.Geo
		if err := tx.Where("code=?", "E01000002").First(&g).Error; err != nil {
			t.Fatalf(err.Error())
		}

		// are we getting a geom.T from a GORM read for .Geometry?
		if g.Geometry.SRID() != 4326 {
			t.Fail()
		}

		// are we getting a geom.T from a GORM read for .LongLatGeom?
		if !reflect.DeepEqual(g.LongLatGeom.FlatCoords(), []float64{-0.09197, 51.51868}) {
			t.Fail()
		}

		// write to database
		ng := model.Geo{
			Name:        "dummy2",
			Code:        "dummy2",
			Geometry:    g.Geometry,
			LongLatGeom: g.LongLatGeom,
			TypeID:      g.TypeID,
			Valid:       false,
		}

		if err := tx.Save(&ng).Error; err != nil {
			t.Fatalf(err.Error())
		}

		// explicit read
		{
			var res model.Geo

			if err := tx.Where("code=?", "dummy2").First(&res).Error; err != nil {
				t.Fatalf(err.Error())
			}

			// are we getting a geom.T from a GORM write/read?
			if res.Geometry.SRID() != 4326 { // XXX
				t.Fail()
			}
			// are we getting a geom.T from a GORM read for .LongLatGeom?
			if !reflect.DeepEqual(res.LongLatGeom.FlatCoords(), []float64{-0.09197, 51.51868}) {
				t.Fail()
			}
		}

		{
			// write nils
			var res model.Geo
			if err := tx.Save(&model.Geo{
				Name:   "nuldummy",
				Code:   "nuldummy",
				TypeID: 6,
				Valid:  false,
			}).Error; err != nil {
				t.Fatalf(err.Error())
			}

			// read nils
			if err := tx.Where("code=?", "nuldummy").First(&res).Error; err != nil {
				t.Fatalf(err.Error())
			}

			if res.Wkbgeometry.Valid {
				t.Fail()
			}

			if res.Wkbgeometry.Valid {
				t.Fail()
			}

		}
	}()

}
