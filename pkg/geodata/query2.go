package geodata

import (
	"context"
	"fmt"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/timer"
)

// Proposed replacement for Query.
// This version separates selecting geocodes from selecting metrics.
func (app *Geodata) Query2(ctx context.Context, year int, bbox, location string, radius int, polygon string, geotypes, geos []string) ([]string, error) {
	err := validateCensusQuery(
		CensusQuerySQLArgs{
			Year:     year,
			Geos:     geos,
			BBox:     bbox,
			Location: location,
			Radius:   radius,
			Polygon:  polygon,
			Geotypes: geotypes,
		},
	)
	if err != nil {
		return nil, err
	}

	sql, err := geocodesSQL(year, bbox, location, radius, polygon, geotypes, geos)
	if err != nil {
		return nil, err
	}

	t := timer.New("query")
	t.Start()
	rows, err := app.db.DB().QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	t.Stop()
	t.Print()
	defer rows.Close()

	var result []string
	tnext := timer.New("next")
	tscan := timer.New("scan")
	for {
		tnext.Start()
		ok := rows.Next()
		tnext.Stop()
		if !ok {
			break
		}

		var geo string

		tscan.Start()
		err := rows.Scan(&geo)
		tscan.Stop()
		if err != nil {
			return nil, err
		}

		result = append(result, geo)
	}
	tnext.Print()
	tscan.Print()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func geocodesSQL(year int, bbox, location string, radius int, polygon string, geotypes, geos []string) (string, error) {
	var geoConditions string
	if !wantAllRows(geos) {
		// fetch conditions SQL
		geoCondition, geoErr := geoSQL(geos)
		bboxCondition, bboxErr := bboxSQL(bbox)
		radiusCondition, radiusErr := radiusSQL(location, radius)
		polygonCondition, polygonErr := polygonSQL(polygon)

		// check errs, return on first found
		for _, err := range []error{
			geoErr,
			bboxErr,
			radiusErr,
			polygonErr,
		} {
			if err != nil {
				return "", err
			}
		}

		// collate join conditions with sql OR
		var conditions []string
		for _, condition := range []string{
			geoCondition,
			bboxCondition,
			radiusCondition,
			polygonCondition,
		} {
			if condition != "" {
				conditions = append(conditions, condition)
			}
		}
		geoConditions = fmt.Sprintf(
			"AND (\n    %s)\n",
			strings.Join(conditions, "    OR\n"),
		)
	}

	// construct WHERE condition for geotypes
	geotypeConditions, err := geotypeSQL("geo_type.name", geotypes)
	if err != nil {
		return "", err
	}

	// construct SQL
	template := `
SELECT
	geo.code AS geography_code
FROM
	geo,
	geo_type
WHERE geo.valid
AND geo_type.id = geo.type_id
	-- geotype conditions:
%s
	-- geo conditions:
%s
`

	sql := fmt.Sprintf(
		template,
		geotypeConditions,
		geoConditions,
	)
	return sql, nil
}
