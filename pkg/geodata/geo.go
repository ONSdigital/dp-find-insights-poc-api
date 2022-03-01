package geodata

import (
	"context"
	"database/sql"
	"encoding/json"
)

type GeoJSON struct {
	Type     string            `json:"type"`
	Features []json.RawMessage `json:"features"`
}

func (app *Geodata) Geo(ctx context.Context, year int, region string) (*GeoJSON, error) {
	template := `
		SELECT
			ST_AsGeoJSON(wkb_long_lat_geom),
			ST_AsGeoJSON(wkb_geometry),
			ST_ASGeoJSON(ST_BoundingDiagonal(wkb_geometry))
		FROM geo
		WHERE code = $1 
		`

	stmt, err := app.db.DB().PrepareContext(ctx, template)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var centroid, boundary, bbox string
	err = stmt.QueryRowContext(ctx, region).Scan(&centroid, &boundary, &bbox)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoContent
		}
		return nil, err
	}

	return &GeoJSON{
		Type: "FeatureCollection",
		Features: []json.RawMessage{
			json.RawMessage(centroid),
			json.RawMessage(boundary),
			json.RawMessage(bbox),
		},
	}, nil
}
