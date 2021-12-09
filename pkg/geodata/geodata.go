package geodata

import (
	"context"
	"fmt"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/table"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/timer"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/where"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Geodata struct {
	db         *database.Database
	maxMetrics int
}

func New(db *database.Database, maxMetrics int) (*Geodata, error) {
	return &Geodata{
		db:         db,
		maxMetrics: maxMetrics,
	}, nil
}

func (app *Geodata) Query(ctx context.Context, dataset, bbox, location string, radius int, polygon string, geotypes, rows, cols []string) (string, error) {
	if len(bbox) > 0 {
		return app.bboxQuery(ctx, bbox, geotypes, cols)
	}
	if len(location) > 0 {
		return app.radiusQuery(ctx, location, radius, geotypes, cols)
	}
	if len(polygon) > 0 {
		return app.polygonQuery(ctx, polygon, geotypes, cols)
	}
	return app.rowQuery(ctx, rows, cols)
}

// rowQuery returns the csv table for the given geometry and category codes.
//
func (app *Geodata) rowQuery(ctx context.Context, geos, cats []string) (string, error) {

	if len(geos) == 0 && len(cats) == 0 {
		return "", ErrMissingParams
	}

	// Construct SQL
	//
	template := `
SELECT
    geo.code AS geography_code,
    nomis_category.long_nomis_code AS category_code,
    geo_metric.metric AS value
FROM
    geo_metric,
    geo,
    nomis_category,
    data_ver
WHERE data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
%s
AND geo_metric.data_ver_id = data_ver.id
AND geo_metric.geo_id = geo.id
%s
AND nomis_category.year = %d
AND geo_metric.category_id = nomis_category.id
`

	geoWhere, err := additionalCondition("geo.code", geos)
	if err != nil {
		return "", err
	}

	catWhere, err := additionalCondition("nomis_category.long_nomis_code", cats)
	if err != nil {
		return "", err
	}

	sql := fmt.Sprintf(
		template,
		geoWhere,
		catWhere,
		2011,
	)
	fmt.Printf("sql: %s\n", sql)

	return app.collectCells(ctx, sql)
}

// bboxQuery returns the csv table for areas intersecting with the given bbox
//
func (app *Geodata) bboxQuery(ctx context.Context, bbox string, geotypes, cats []string) (string, error) {
	var p1lon, p1lat, p2lon, p2lat float64
	fields, err := fmt.Sscanf(bbox, "%f,%f,%f,%f", &p1lon, &p1lat, &p2lon, &p2lat)
	if err != nil {
		return "", fmt.Errorf("scanning bbox %q: %w", bbox, err)
	}
	if fields != 4 {
		return "", fmt.Errorf("bbox missing a number: %w", ErrMissingParams)
	}

	// Construct SQL
	//
	template := `
SELECT
	geo.code AS geography_code,
	nomis_category.long_nomis_code AS category_code,
	geo_metric.metric AS value
FROM
	geo,
	geo_type,
	geo_metric,
	data_ver,
	nomis_category
WHERE geo.wkb_geometry && ST_GeomFromText(
		'MULTIPOINT(%f %f, %f %f)',
		4326
	)
AND geo.valid
AND geo.type_id = geo_type.id
%s
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = %d
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = %d
%s
`

	catWhere, err := additionalCondition("nomis_category.long_nomis_code", cats)
	if err != nil {
		return "", err
	}

	geotypeWhere, err := additionalCondition("geo_type.name", geotypes)
	if err != nil {
		return "", err
	}

	sql := fmt.Sprintf(
		template,
		p1lon,
		p1lat,
		p2lon,
		p2lat,
		geotypeWhere,
		2011,
		2011,
		catWhere,
	)

	fmt.Printf("sql: %s\n", sql)

	return app.collectCells(ctx, sql)
}

// radiusQuery returns the csv table for areas within radius meters from location
//
func (app *Geodata) radiusQuery(ctx context.Context, location string, radius int, geotypes, cats []string) (string, error) {
	var lon, lat float64
	fields, err := fmt.Sscanf(location, "%f,%f", &lon, &lat)
	if err != nil {
		return "", fmt.Errorf("scanning location %q: %w", location, err)
	}
	if fields != 2 {
		return "", fmt.Errorf("location missing a number: %w", ErrMissingParams)
	}

	if radius < 1 {
		return "", fmt.Errorf("radius must be >0: %q: %w", radius, err)
	}

	template := `
SELECT
    geo.code AS geography_code,
    nomis_category.long_nomis_code AS category_code,
    geo_metric.metric AS value
FROM
    geo,
    geo_type,
    geo_metric,
    data_ver,
    nomis_category
WHERE ST_DWithin(
    geo.wkb_long_lat_geom::geography,
    ST_SetSRID(
        ST_Point(%f, %f),
        4326
    )::geography,
    %d
)
AND geo.valid
AND geo_type.id = geo.type_id
%s
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = %d
ANd data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = %d
%s
`

	catWhere, err := additionalCondition("nomis_category.long_nomis_code", cats)
	if err != nil {
		return "", err
	}

	geotypeWhere, err := additionalCondition("geo_type.name", geotypes)
	if err != nil {
		return "", err
	}

	sql := fmt.Sprintf(
		template,
		lon,
		lat,
		radius,
		geotypeWhere,
		2011,
		2011,
		catWhere,
	)

	fmt.Printf("sql: %s\n", sql)

	return app.collectCells(ctx, sql)
}

// polygonQuery returns the csv table for areas within radius meters from location
//
func (app *Geodata) polygonQuery(ctx context.Context, polygon string, geotypes, cats []string) (string, error) {
	points, err := ParsePolygon(polygon)
	if err != nil {
		return "", fmt.Errorf("parsing polygon: %q: %w", polygon, err)
	}

	template := `
SELECT
    geo.code AS geography_code,
    nomis_category.long_nomis_code AS category_code,
    geo_metric.metric AS value
FROM
    geo,
    geo_type,
    geo_metric,
    data_ver,
    nomis_category
WHERE ST_Covers(
    ST_Polygon(
        'LINESTRING (%s)'::geometry,
        4326
    ),
    geo.wkb_geometry
)
AND geo.valid
AND geo_type.id = geo.type_id
%s
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = %d
ANd data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = %d
%s
`

	// convert the slice of Points to a slice of strings holding coordinates like "lon lat"
	var coords []string
	for _, point := range points {
		coords = append(coords, point.String())
	}

	// join the coordinate strings into a form usable in LINESTRING(...)
	linestring := strings.Join(coords, ",")

	catWhere, err := additionalCondition("nomis_category.long_nomis_code", cats)
	if err != nil {
		return "", err
	}

	geotypeWhere, err := additionalCondition("geo_type.name", geotypes)
	if err != nil {
		return "", err
	}

	sql := fmt.Sprintf(
		template,
		linestring,
		geotypeWhere,
		2011,
		2011,
		catWhere,
	)

	fmt.Printf("sql: %s\n", sql)

	return app.collectCells(ctx, sql)
}

// collectCells runs the query in sql and returns the results as a csv.
// sql must be a query against the geo_metric table selecting exactly
// code, category and metric.
//
func (app *Geodata) collectCells(ctx context.Context, sql string) (string, error) {
	// Allocate output table
	//
	tbl, err := table.New("geography_code")
	if err != nil {
		return "", err
	}

	// Set up output buffer
	//
	var body strings.Builder
	body.Grow(1000000)

	// Query the db.
	//
	t := timer.New("query")
	t.Start()
	rows, err := app.db.DB().QueryContext(ctx, sql)
	if err != nil {
		return "", err
	}
	t.Stop()
	t.Print()
	defer rows.Close()

	tnext := timer.New("next")
	tscan := timer.New("scan")
	var nmetrics int
	for {
		tnext.Start()
		ok := rows.Next()
		tnext.Stop()
		if !ok {
			break
		}

		if app.maxMetrics > 0 {
			nmetrics++
			if nmetrics > app.maxMetrics {
				return "", ErrTooManyMetrics
			}
		}

		var geo string
		var cat string
		var value float64

		tscan.Start()
		err := rows.Scan(&geo, &cat, &value)
		tscan.Stop()
		if err != nil {
			return "", err
		}

		tbl.SetCell(geo, cat, value)
	}
	tnext.Print()
	tscan.Print()

	if err := rows.Err(); err != nil {
		return "", err
	}

	tgen := timer.New("generate")
	tgen.Start()
	err = tbl.Generate(&body)
	tgen.Stop()
	tgen.Print()
	if err != nil {
		return "", err
	}

	return body.String(), nil
}

// additionalCondition wraps the output of WherePart inside "AND (...)".
// We "know" this additionalCondition will not be the first additionalCondition in the query.
func additionalCondition(col string, args []string) (string, error) {
	if len(args) == 0 {
		return "", nil
	}
	body, err := where.WherePart(col, args)
	if err != nil {
		return "", err
	}

	template := `
AND (
%s
)`
	return fmt.Sprintf(template, body), nil
}
