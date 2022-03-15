package geodata

import (
	"context"
	"database/sql"

	"github.com/ONSdigital/dp-find-insights-poc-api/sentinel"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkb"
)

func (app *Geodata) Geo(ctx context.Context, year int, geocode string) (*geojson.FeatureCollection, error) {
	template := `
		SELECT
			ST_AsBinary(wkb_long_lat_geom),
			ST_AsBinary(wkb_geometry),
			ST_AsBinary(ST_BoundingDiagonal(wkb_geometry))
		FROM geo
		WHERE code = $1 
		`

	stmt, err := app.db.DB().PrepareContext(ctx, template)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var centroid, boundary, bbox []byte
	err = stmt.QueryRowContext(ctx, geocode).Scan(&centroid, &boundary, &bbox)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sentinel.ErrNoContent
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

	return collection, nil
}
