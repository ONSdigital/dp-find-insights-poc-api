package geobb

import (
	"bytes"
	"encoding/json"
	"log"
	"math"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"gorm.io/gorm"
)

type GeoBB struct {
	Gdb *gorm.DB
}

type GeoObj struct {
	En      string    `json:"en"`
	Cy      string    `json:"cy"`
	GeoType string    `json:"geoType"`
	GeoCode string    `json:"geoCode"`
	Bbox    []float64 `json:"bbox"`
}
type GeoObjs []GeoObj

type Params struct {
	Pretty bool
	Welsh  bool
	Geos   []string
}

var geotype = map[string]int{
	"LAD":  4,
	"MSOA": 5,
}

func (g *GeoBB) AsJSON(params Params) string {
	var geos []model.Geo
	if err := g.Gdb.Order("code").Where("type_id in (?,?)", geotype[params.Geos[0]], geotype[params.Geos[1]]).Find(&geos).Error; err != nil {
		log.Print(err)
	}

	var lads GeoObjs

	for _, geo := range geos {

		wkb := geo.Wkbgeometry

		// null
		if !wkb.Valid {
			continue
		}

		// as gorm create hook?
		geomt, err := ewkbhex.Decode(wkb.String)
		if err != nil {
			log.Print(err)
		}

		// x = long/0 y = lat/1

		minx := roundToFive(geomt.Bounds().Min(0))
		miny := roundToFive(geomt.Bounds().Min(1))
		maxx := roundToFive(geomt.Bounds().Max(0))
		maxy := roundToFive(geomt.Bounds().Max(1))

		// should really be join
		m := make(map[int32]string)
		m[4] = "LAD"
		m[5] = "MSOA"

		var welshName string
		if params.Welsh {
			welshName = geo.WelshName
		}

		lad := GeoObj{
			En:      geo.Name,
			GeoType: m[geo.TypeID],
			GeoCode: geo.Code,
			Cy:      welshName,
			Bbox:    []float64{minx, miny, maxx, maxy},
		}

		lads = append(lads, lad)

	}

	b, err := jsonMarshallNoEsc(&lads)
	if err != nil {
		log.Print(err)
	}

	if params.Pretty {
		var out bytes.Buffer
		if err := json.Indent(&out, b, " ", " "); err != nil {
			log.Print(err)
		}

		return out.String()
	}

	return string(b)
}

func roundToFive(f float64) float64 {
	return math.Round((f * 100000)) / 100000
}

func jsonMarshallNoEsc(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(&v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
