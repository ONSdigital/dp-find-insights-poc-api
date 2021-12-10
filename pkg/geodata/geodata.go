package geodata

import (
	"context"
	"errors"
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
	if dataset == "census" {
		return app.censusQuery(ctx, rows, bbox, location, radius, polygon, geotypes, cols)
	}

	// geography_code is hardcoded for compatibility with previous skinny queries
	compatInclude := []string{"geography_code"}
	if len(bbox) > 0 {
		return app.bboxQuery(ctx, bbox, geotypes, cols, compatInclude)
	}
	if len(location) > 0 {
		return app.radiusQuery(ctx, location, radius, geotypes, cols, compatInclude)
	}
	if len(polygon) > 0 {
		return app.polygonQuery(ctx, polygon, geotypes, cols, compatInclude)
	}
	return app.rowQuery(ctx, rows, geotypes, cols, compatInclude)
}

// rowQuery returns the csv table for the given geometry and category codes.
//
func (app *Geodata) rowQuery(ctx context.Context, geos, geotypes, cats, include []string) (string, error) {

	if len(geos) == 0 && len(cats) == 0 {
		return "", ErrMissingParams
	}

	// Construct SQL
	//
	template := `
SELECT
    geo.code AS geography_code,
    geo_type.name AS geotype,
    nomis_category.long_nomis_code AS category_code,
    geo_metric.metric AS value
FROM
    geo_metric,
    geo_type,
    geo,
    nomis_category,
    data_ver
WHERE data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
%s
AND geo_metric.data_ver_id = data_ver.id
AND geo_metric.geo_id = geo.id
%s
AND geo.valid
AND geo.type_id = geo_type.id
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

	geotypeWhere, err := additionalCondition("geo_type.name", geotypes)
	if err != nil {
		return "", err
	}

	sql := fmt.Sprintf(
		template,
		geoWhere,
		catWhere,
		geotypeWhere,
		2011,
	)
	fmt.Printf("sql: %s\n", sql)

	return app.collectCells(ctx, sql, include)
}

// bboxQuery returns the csv table for areas intersecting with the given bbox
//
func (app *Geodata) bboxQuery(ctx context.Context, bbox string, geotypes, cats, include []string) (string, error) {
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
    geo_type.name AS geotype,
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

	return app.collectCells(ctx, sql, include)
}

// radiusQuery returns the csv table for areas within radius meters from location
//
func (app *Geodata) radiusQuery(ctx context.Context, location string, radius int, geotypes, cats, include []string) (string, error) {
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
    geo_type.name AS geotype,
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
		lon,
		lat,
		radius,
		geotypeWhere,
		2011,
		2011,
		catWhere,
	)

	fmt.Printf("sql: %s\n", sql)

	return app.collectCells(ctx, sql, include)
}

// polygonQuery returns the csv table for areas within radius meters from location
//
func (app *Geodata) polygonQuery(ctx context.Context, polygon string, geotypes, cats, include []string) (string, error) {
	points, err := ParsePolygon(polygon)
	if err != nil {
		return "", fmt.Errorf("parsing polygon: %q: %w", polygon, err)
	}

	template := `
SELECT
    geo.code AS geography_code,
    geo_type.name AS geotype,
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
AND data_ver.ver_string = '2.2'
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

	return app.collectCells(ctx, sql, include)
}

// collectCells runs the query in sql and returns the results as a csv.
// sql must be a query against the geo_metric table selecting exactly
// code, category and metric.
//
func (app *Geodata) collectCells(ctx context.Context, sql string, include []string) (string, error) {
	// Allocate output table
	//
	tbl := table.New()

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
		var geotype string
		var cat string
		var value float64

		tscan.Start()
		err := rows.Scan(&geo, &geotype, &cat, &value)
		tscan.Stop()
		if err != nil {
			return "", err
		}

		tbl.SetCell(geo, geotype, cat, value)
	}
	tnext.Print()
	tscan.Print()

	if err := rows.Err(); err != nil {
		return "", err
	}

	tgen := timer.New("generate")
	tgen.Start()
	err = tbl.Generate(&body, include)
	tgen.Stop()
	tgen.Print()
	if err != nil {
		return "", err
	}

	return body.String(), nil
}

// censusQuery is the merged query which is the logical OR of the other specific queries.
// The specific "skinny" queries can probably go away soon.
//
// Any combination of rows, bbox, radius and polygon queries can be given, and geotype can be exposed in the csv output.
//
// Although this query method is not complicated, it is too long.
// Break it up in the fullness of time.
//
func (app *Geodata) censusQuery(ctx context.Context, geos []string, bbox, location string, radius int, polygon string, geotypes, cols []string) (string, error) {
	var conditions []string

	// construct conditions for explicitly named rows=
	if len(geos) > 0 {
		condition, err := where.WherePart("geo.code", geos)
		if err != nil {
			return "", err
		}
		conditions = append(conditions, condition)
	}

	// construct condition for bbox=
	if bbox != "" {
		var p1lon, p1lat, p2lon, p2lat float64
		fields, err := fmt.Sscanf(bbox, "%f,%f,%f,%f", &p1lon, &p1lat, &p2lon, &p2lat)
		if err != nil {
			return "", fmt.Errorf("scanning bbox %q: %w", bbox, err)
		}
		if fields != 4 {
			return "", fmt.Errorf("bbox missing a number: %w", ErrMissingParams)
		}
		condition := fmt.Sprintf(`
    geo.wkb_geometry && ST_GeomFromText(
        'MULTIPOINT(%f %f, %f %f)',
        4326
    )
`,
			p1lon,
			p1lat,
			p2lon,
			p2lat,
		)
		conditions = append(conditions, condition)
	}

	// construct condition for location= and radius=
	if location != "" || radius > 0 {
		if location == "" || radius == 0 {
			return "", fmt.Errorf("radius queries require both location (%s) and radius (%d)", location, radius)
		}
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
		condition := fmt.Sprintf(`
    ST_DWithin(
        geo.wkb_long_lat_geom::geography,
        ST_SetSRID(
            ST_Point(%f, %f),
            4326
        )::geography,
        %d
    )`,
			lon,
			lat,
			radius,
		)
		conditions = append(conditions, condition)
	}

	// construct condition for polygon=
	if polygon != "" {
		points, err := ParsePolygon(polygon)
		if err != nil {
			return "", fmt.Errorf("parsing polygon: %q: %w", polygon, err)
		}
		// convert the slice of Points to a slice of strings holding coordinates like "lon lat"
		var coords []string
		for _, point := range points {
			coords = append(coords, point.String())
		}
		// join the coordinate strings into a form usable in LINESTRING(...)
		linestring := strings.Join(coords, ",")
		condition := fmt.Sprintf(`
    ST_COVERS(
        ST_Polygon(
            'LINESTRING (%s)'::geometry,
            4326
        ),
        geo.wkb_geometry
    )`,
			linestring,
		)
		conditions = append(conditions, condition)
	}

	if len(conditions) == 0 {
		return "", errors.New("must specify a condition (rows,bbox,location/radius, or polygon)")
	}

	// join conditions with sql OR
	geoConditions := strings.Join(conditions, "    OR\n")

	// construct WHERE condition for geotypes
	geotypeConditions, err := additionalCondition("geo_type.name", geotypes)
	if err != nil {
		return "", err
	}

	// split column list into includes and categories
	include, cats := splitCols(cols)

	// construct WHERE condition for categories
	catConditions, err := additionalCondition("nomis_category.long_nomis_code", cats)
	if err != nil {
		return "", err
	}

	// construct final SQL
	template := `
SELECT
    geo.code AS geography_code,
    geo_type.name AS geotype,
    nomis_category.long_nomis_code AS category_code,
    geo_metric.metric AS value
FROM
    geo,
    geo_type,
    geo_metric,
    data_ver,
    nomis_category
WHERE (
    -- geo conditions:
%s
)
AND geo.valid
AND geo_type.id = geo.type_id
    -- geotype conditions:
%s
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = 2011
    -- category conditions:
%s
`
	sql := fmt.Sprintf(
		template,
		geoConditions,
		geotypeConditions,
		catConditions,
	)

	fmt.Printf("sql: %s\n", sql)

	return app.collectCells(ctx, sql, include)
}

// splitCols separates special column names from geography names
// (XXXX This duplicates some of the work in where.GetValues().
// Think of a better way.)
func splitCols(cols []string) (include []string, cats []string) {
	for _, instance := range cols {
		tokens := strings.Split(instance, ",")
		for _, token := range tokens {
			if token == table.ColGeographyCode || token == table.ColGeotype {
				include = append(include, token)
			} else {
				cats = append(cats, token)
			}
		}
	}
	return
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
