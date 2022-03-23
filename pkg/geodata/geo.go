package geodata

import (
	"context"
	"database/sql"

	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkb"
)

type Resp struct {
	Meta struct {
		Name    string `json:"name"`
		Code    string `json:"code"`
		Geotype string `json:"geotype"`
	} `json:"meta"`
	GeoJSON *geojson.FeatureCollection `json:"geo_json"`
}

func (app *Geodata) Geo(ctx context.Context, year int, geocode string, geoname string) (*Resp, error) {
	queryString := `
	SELECT
		ST_AsBinary(geo.wkb_long_lat_geom),
		ST_AsBinary(geo.wkb_geometry),
		ST_AsBinary(ST_BoundingDiagonal(geo.wkb_geometry)),
		geo.name,
		geo.code,
		geo_type.name
	FROM
		geo,
		geo_type
	`
	var conditionString string
	var queryCondition string
	if geocode != "" {
		conditionString = `
			WHERE geo.code = $1
			AND geo.type_id = geo_type.id
			`
		queryCondition = geocode
	} else {
		conditionString = `
			WHERE geo.name = $1
			AND geo.type_id = geo_type.id
			`
		queryCondition = geoname
	}
	fullQueryString := queryString + conditionString
	stmt, err := app.db.DB().PrepareContext(ctx, fullQueryString)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var centroid, boundary, bbox []byte
	var name, code, geotype string
	err = stmt.QueryRowContext(ctx, queryCondition).Scan(&centroid, &boundary, &bbox, &name, &code, &geotype)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Resp{}, nil
		}
		return nil, err
	}

	// return no data if there is no geometry (this is the case for England and Wales, Regions, and other geotypes)
	if centroid == nil || boundary == nil || bbox == nil {
		return &Resp{}, nil
	}

	geomCentroid, err := wkb.Unmarshal(centroid)
	if err != nil {
		return nil, err
	}
	geomBoundary, err := wkb.Unmarshal(boundary)
	if err != nil {
		return nil, err
	}
	geomBbox, err := wkb.Unmarshal(bbox)
	if err != nil {
		return nil, err
	}

	collection := &geojson.FeatureCollection{
		Features: []*geojson.Feature{
			{ID: "centroid", Geometry: geomCentroid},
			{ID: "boundary", Geometry: geomBoundary},
			{ID: "bbox", Geometry: geomBbox},
		},
	}

	r := &Resp{GeoJSON: collection}
	r.Meta.Name = name
	r.Meta.Code = code
	r.Meta.Geotype = geotype

	return r, nil
}
