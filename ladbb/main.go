package ladbb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
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

type LadGeom struct {
	Gdb *gorm.DB
}

func (g *LadGeom) AsJSON() string {
	var geos []model.Geo
	if err := g.Gdb.Where("type_id=4").Find(&geos).Error; err != nil {
		log.Print(err)
	}

	j := make(LadBB)
	for _, geo := range geos {
		geomt := geo.Geometry

		if geomt == nil {
			continue
		}

		x, y := GetXYOrder(geomt.Bounds().Layout().String())

		f := "%.5f"

		j[geo.Code] = LadBBValues{Name: geo.Name,
			Lat:  fmt.Sprintf(f, geo.Lat),
			Lon:  fmt.Sprintf(f, geo.Long),
			MaxX: fmt.Sprintf(f, geomt.Bounds().Max(x)),
			MaxY: fmt.Sprintf(f, geomt.Bounds().Max(y)),
			MinX: fmt.Sprintf(f, geomt.Bounds().Min(x)),
			MinY: fmt.Sprintf(f, geomt.Bounds().Min(y)),
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

	return out.String()
}

func GetXYOrder(xy string) (x, y int) {

	if xy == "XY" {
		return 0, 1
	}

	return 1, 0
}
