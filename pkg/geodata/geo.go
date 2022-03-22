package geodata

import (
	"context"
	"database/sql"

	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkb"
)

type Resp struct {
	Meta struct {
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"meta"`
	GeoJSON *geojson.FeatureCollection `json:"geo_json"`
}

func (app *Geodata) Geo(ctx context.Context, year int, geocode string) (*Resp, error) {
	template := `
		SELECT
			ST_AsBinary(wkb_long_lat_geom),
			ST_AsBinary(wkb_geometry),
			ST_AsBinary(ST_BoundingDiagonal(wkb_geometry)),
			name,
			code
		FROM geo
		WHERE code = $1 
		`

	stmt, err := app.db.DB().PrepareContext(ctx, template)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var centroid, boundary, bbox []byte
	var name, code string
	err = stmt.QueryRowContext(ctx, geocode).Scan(&centroid, &boundary, &bbox, &name, &code)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Resp{}, nil
		}
		return nil, err
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

	return r, nil
}
