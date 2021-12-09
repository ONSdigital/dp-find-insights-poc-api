package main

import (
	"fmt"
	"log"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/spf13/cast"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// populates geo.wkb_long_lat_geom with long, lat POINT
func main() {
	db, err := gorm.Open(postgres.Open(database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var geos []model.Geo

	db.Find(&geos)

	for i := range geos {
		g := geos[i]
		if g.Lat != 0 && g.Long != 0 && g.Valid {
			fmt.Printf("%#v\n", g.ID)
			update(db, g.ID, g.Long, g.Lat)
		}
	}
}

func update(db *gorm.DB, id int32, long, lat float64) {
	// different SRID syntax?
	if err := db.Exec("UPDATE geo SET wkb_long_lat_geom=( SELECT ST_GeomFromText('SRID=4326;POINT('|| ? || ' ' || ? || ')') ) WHERE id=?", cast.ToString(long), cast.ToString(lat), id).Error; err != nil {
		log.Print(err)
	}
}
