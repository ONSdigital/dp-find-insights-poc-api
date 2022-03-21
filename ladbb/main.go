package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/spf13/cast"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LadBBValues struct {
	Name string
	Lat  string
	Lon  string
	MinX string
	MinY string
	MaxX string
	MaxY string
}

type LadBB map[string]LadBBValues

func main() {

	gdb, err := gorm.Open(postgres.Open(database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Print(err)
	}

	var geos []model.Geo
	gdb.Where("type_id=4").Find(&geos)

	j := make(LadBB)
	for _, geo := range geos {

		// as gorm create hook?
		geomt, err := ewkbhex.Decode(geo.Wkb_geometry)
		if err != nil {
			log.Print(err)
		}

		if geomt.Bounds().Layout().String() == "XY" {
			j[geo.Code] = LadBBValues{Name: geo.Name,
				Lat:  cast.ToString(geo.Lat),
				Lon:  cast.ToString(geo.Long),
				MaxX: cast.ToString(geomt.Bounds().Max(0)),
				MaxY: cast.ToString(geomt.Bounds().Max(1)),
				MinX: cast.ToString(geomt.Bounds().Min(0)),
				MinY: cast.ToString(geomt.Bounds().Min(1)),
			}
		} else {
			log.Fatal("unsupported layout")
		}

	}

	b, err := json.Marshal(&j)
	if err != nil {
		log.Print(err)
	}

	var out bytes.Buffer
	if err := json.Indent(&out, b, " ", " "); err != nil {
		log.Print(err)
	}

	fmt.Println(out.String())
}
