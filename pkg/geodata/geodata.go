package geodata

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

func (app *Geodata) Query(ctx context.Context, dataset, bbox, location string, radius int, polygon string, geotypes, rows, cols []string, censustable string) (string, error) {
	if dataset == "census" {
		return app.censusQuery(ctx, rows, bbox, location, radius, polygon, geotypes, cols, censustable)
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

	if nmetrics == 0 {
		return "", ErrNoContent
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

type CensusQuerySQLArgs struct {
	Geos        []string
	BBox        string
	Location    string
	Radius      int
	Polygon     string
	Geotypes    []string
	Cols        []string
	Censustable string
}

// censusQuery is the merged query which is the logical OR of the other specific queries.
// The specific "skinny" queries can probably go away soon.
//
// Any combination of rows, bbox, radius and polygon queries can be given, and geotype can be exposed in the csv output.
//
// Although this query method is not complicated, it is too long.
// Break it up in the fullness of time.
//
func (app *Geodata) censusQuery(ctx context.Context, geos []string, bbox, location string, radius int, polygon string, geotypes, cols []string, censustable string) (string, error) {

	sql, include, err := CensusQuerySQL(
		ctx,
		CensusQuerySQLArgs{
			Geos:        geos,
			BBox:        bbox,
			Location:    location,
			Radius:      radius,
			Polygon:     polygon,
			Geotypes:    geotypes,
			Cols:        cols,
			Censustable: censustable,
		},
	)
	if err != nil {
		return "", err
	}

	fmt.Printf("sql: %s\n", sql)

	return app.collectCells(ctx, sql, include)
}

func CensusQuerySQL(ctx context.Context, args CensusQuerySQLArgs) (sql string, include []string, err error) {
	// validate args
	if err := validateCensusQuery(args); err != nil {
		return sql, include, err
	}

	// fetch conditions SQL
	geoCondition, geoErr := geoSQL(args.Geos)
	bboxCondition, bboxErr := bboxSQL(args.BBox)
	radiusCondition, radiusErr := radiusSQL(args.Location, args.Radius)
	polygonCondition, polygonErr := polygonSQL(args.Polygon)

	// check errs, return on first found
	for _, err := range []error{
		geoErr,
		bboxErr,
		radiusErr,
		polygonErr,
	} {
		if err != nil {
			return sql, include, err
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
	geoConditions := strings.Join(conditions, "    OR\n")

	// construct WHERE condition for geotypes
	geotypeConditions, err := additionalCondition("geo_type.name", args.Geotypes)
	if err != nil {
		return sql, include, err
	}

	// split column list into includes and categories
	include, cats := splitCols(args.Cols)

	// construct WHERE condition for categories
	catConditions, err := categorySQL(cats, args.Censustable)
	if err != nil {
		return sql, include, err
	}

	// construct additional conditions for censustable / short_nomis_code
	censustableFromSQL, censustableAndSQL := censusTableFromAndSQL(args.Censustable)

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
	%s
WHERE (
    -- geo conditions:
	%s
)
AND geo.valid
AND geo_type.id = geo.type_id
    -- geotype conditions:
%s
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
	sql = fmt.Sprintf(
		template,
		censustableFromSQL,
		geoConditions,
		geotypeConditions,
		censustableAndSQL,
		catConditions,
	)
	return sql, include, nil
}

func validateCensusQuery(args CensusQuerySQLArgs) error {
	// 'conditions' cant all be null / default
	if len(args.Geos) == 0 &&
		args.BBox == "" &&
		args.Location == "" &&
		args.Radius == 0 &&
		args.Polygon == "" {
		return errors.New("must specify a condition (rows, bbox, location/radius, and/or polygon)")
	}
	return nil
}

func geoSQL(geos []string) (string, error) {
	if len(geos) > 0 {
		return where.WherePart("geo.code", geos)
	}
	return "", nil
}

func bboxSQL(bbox string) (string, error) {
	if bbox != "" {
		coords := []float64{}
		for _, coordStr := range strings.Split(bbox, ",") {
			coord, err := strconv.ParseFloat(coordStr, 64)
			if err != nil {
				return "", fmt.Errorf("error parsing bbox %q: %w", bbox, err)
			}
			coords = append(coords, coord)
		}
		if len(coords) != 4 {
			return "", fmt.Errorf("valid bbox is 'lon,lat,lon,lat', received %q: %w", bbox, ErrInvalidParams)
		}
		sql := fmt.Sprintf(`
geo.wkb_geometry && ST_GeomFromText(
	'MULTIPOINT(%f %f, %f %f)',
	4326
)
`,
			coords[0],
			coords[1],
			coords[2],
			coords[3],
		)
		return sql, nil
	}
	return "", nil
}

func radiusSQL(location string, radius int) (string, error) {
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
		sql := fmt.Sprintf(`
ST_DWithin(
	geo.wkb_long_lat_geom::geography,
	ST_SetSRID(
		ST_Point(%f, %f),
		4326
	)::geography,
	%d
)
`,
			lon,
			lat,
			radius,
		)
		return sql, nil
	}
	return "", nil
}

func polygonSQL(polygon string) (string, error) {
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
		sql := fmt.Sprintf(`
ST_COVERS(
	ST_Polygon(
		'LINESTRING (%s)'::geometry,
		4326
	),
	geo.wkb_geometry
)
`,
			linestring,
		)
		return sql, nil
	}
	return "", nil
}

func censusTableFromAndSQL(censustable string) (string, string) {
	var fromSQL string
	var andSQL string
	if censustable != "" {
		fromSQL = ", nomis_desc"
		andSQL = fmt.Sprintf(
			`AND nomis_desc.short_nomis_code = '%s'`,
			censustable,
		)
	}
	return fromSQL, andSQL
}

func categorySQL(namedCats []string, censusTable string) (string, error) {
	if len(namedCats) == 0 && censusTable == "" {
		return "", nil
	}
	// get sql for selecting named categories
	var namedCatSQL string
	var err error
	if len(namedCats) > 0 {
		namedCatSQL, err = where.WherePart("nomis_category.long_nomis_code", namedCats)
		if err != nil {
			return "", err
		}
	} else {
		namedCatSQL = ""
	}
	// get sql for selecting categories by nomis desc selection - will need joining 'OR' if there were named cats
	var censusTableSQL string
	if censusTable != "" {
		censusTableSQL = " nomis_category.nomis_desc_id = nomis_desc.id"
		if namedCatSQL != "" {
			censusTableSQL = " OR\n " + censusTableSQL
		}
	}
	template := `
AND (
%s
%s
)`
	return fmt.Sprintf(template, namedCatSQL, censusTableSQL), nil
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
