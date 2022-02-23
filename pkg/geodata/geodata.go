package geodata

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/table"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/timer"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/where"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const allRowsToken = "ALL" // rows= token that means grab all rows, as in rows=ALL

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

func (app *Geodata) Query(ctx context.Context, year int, bbox, location string, radius int, polygon string, geotypes, rows, cols []string, censustable string) (string, error) {
	return app.censusQuery(ctx, year, rows, bbox, location, radius, polygon, geotypes, cols, censustable)
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

		nmetrics++
		if app.maxMetrics > 0 {
			if nmetrics > app.maxMetrics {
				return "", fmt.Errorf("%w: limit is %d", ErrTooManyMetrics, app.maxMetrics)
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
	Year        int
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
func (app *Geodata) censusQuery(ctx context.Context, year int, geos []string, bbox, location string, radius int, polygon string, geotypes, cols []string, censustable string) (string, error) {

	sql, include, err := CensusQuerySQL(
		ctx,
		CensusQuerySQLArgs{
			Year:        year,
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

	var geoConditions string
	if !wantAllRows(args.Geos) {
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
		geoConditions = fmt.Sprintf(
			"AND (\n    %s)\n",
			strings.Join(conditions, "    OR\n"),
		)
	}

	// construct WHERE condition for geotypes
	geotypeConditions, err := geotypeSQL("geo_type.name", args.Geotypes)
	if err != nil {
		return sql, include, err
	}

	// parse cols query strings into a ValueSet
	catset, err := where.ParseMultiArgs(args.Cols)
	if err != nil {
		return sql, include, err
	}

	// extract special column names from ValueSet
	include, catset, err = ExtractSpecialCols(catset)
	if err != nil {
		return sql, include, err
	}

	// construct WHERE condition for categories
	catConditions, err := categorySQL(catset, args.Censustable)
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
WHERE geo.valid
AND geo_type.id = geo.type_id
    -- geotype conditions:
%s
	-- geo conditions:
%s
%s
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = %d
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = data_ver.census_year
    -- category conditions:
%s
`
	sql = fmt.Sprintf(
		template,
		censustableFromSQL,
		geotypeConditions,
		geoConditions,
		censustableAndSQL,
		args.Year,
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
		return fmt.Errorf("%w: must specify a condition (rows, bbox, location/radius, and/or polygon)", ErrMissingParams)
	}

	set, err := where.ParseMultiArgs(args.Geos)
	if err != nil {
		return err
	}
	return ValidateAllToken(set)
}

// ValidateAllToken verifies that if there is an "ALL" specified, it is the only token.
func ValidateAllToken(set *where.ValueSet) error {
	var all, tokens int

	callback := func(single, low, high *string) (*string, *string, *string, error) {
		var err error
		tokens++
		if single != nil {
			if isAll(*single) {
				all++
			}
		} else {
			if isAll(*low) || isAll(*high) {
				err = fmt.Errorf("%w: ALL cannot be part of a range", ErrInvalidParams)
			}
		}
		// we're only interested in the err status
		return nil, nil, nil, err
	}

	_, err := set.Walk(callback)
	if err != nil {
		return err
	}

	if all == 0 {
		return nil
	}
	if all == 1 && tokens == 1 {
		return nil
	}
	return fmt.Errorf("%w: if used, ALL must be first and only rows= token", ErrInvalidParams)
}

// wantAllRows is true if rows=ALL
func wantAllRows(geos []string) bool {
	return len(geos) == 1 && isAll(geos[0])
}

// isAll is true if token is allRowsToken
func isAll(token string) bool {
	return strings.EqualFold(token, allRowsToken)
}

func geoSQL(geos []string) (string, error) {
	set, err := where.ParseMultiArgs(geos)
	if err != nil {
		return "", err
	}
	return where.WherePart("geo.code", set), nil
}

func bboxSQL(bbox string) (string, error) {
	if bbox != "" {
		coords := []float64{}
		for _, coordStr := range strings.Split(bbox, ",") {
			coord, err := strconv.ParseFloat(coordStr, 64)
			if err != nil {
				return "", fmt.Errorf("%w: error parsing bbox %q: %s", ErrInvalidParams, bbox, err)
			}
			coords = append(coords, coord)
		}
		if len(coords) != 4 {
			return "", fmt.Errorf("%w: valid bbox is 'lon,lat,lon,lat', received %q", ErrInvalidParams, bbox)
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
			return "", fmt.Errorf("%w: radius queries require both location (%s) and radius (%d)", ErrInvalidParams, location, radius)
		}
		var lon, lat float64
		fields, err := fmt.Sscanf(location, "%f,%f", &lon, &lat)
		if err != nil {
			return "", fmt.Errorf("%w: scanning location %q: %s", ErrInvalidParams, location, err)
		}
		if fields != 2 {
			return "", fmt.Errorf("%w: location missing a number", ErrMissingParams)
		}
		if radius < 1 {
			return "", fmt.Errorf("%w: radius must be >0: %q", ErrInvalidParams, radius)
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
			return "", fmt.Errorf("%w: parsing polygon: %q", err, polygon)
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

func categorySQL(set *where.ValueSet, censusTable string) (string, error) {
	var conditions []string

	// get sql for selecting named categories
	namedCatSQL := where.WherePart("nomis_category.long_nomis_code", set)
	if namedCatSQL != "" {
		conditions = append(conditions, namedCatSQL)
	}

	// get sql for selecting categories by nomis desc selection
	if censusTable != "" {
		conditions = append(conditions, " nomis_category.nomis_desc_id = nomis_desc.id")
	}

	if len(conditions) == 0 {
		return "", nil
	}

	template := `
AND (
%s
)`
	return fmt.Sprintf(template, strings.Join(conditions, " OR\n ")), nil
}

// ExtractSpecialCols removes special column names like "geography_code" from
// the ValueSet, returning a reduced ValueSet, and the list of special columns
// found.
func ExtractSpecialCols(set *where.ValueSet) ([]string, *where.ValueSet, error) {
	var includes []string

	callback := func(single, low, high *string) (*string, *string, *string, error) {
		var err error
		if single != nil {
			if *single == table.ColGeographyCode || *single == table.ColGeotype {
				includes = append(includes, *single)
				single = nil
			}
		} else {
			if *low == table.ColGeographyCode || *high == table.ColGeographyCode || *low == table.ColGeotype || *high == table.ColGeotype {
				err = fmt.Errorf("%w: special columns cannot be part of a range", ErrInvalidParams)
			}
		}
		if err != nil {
			return nil, nil, nil, err
		}
		return single, low, high, nil
	}

	newset, err := set.Walk(callback)
	return includes, newset, err
}

// geotypeSQL generates an AND where part for geotypes.
func geotypeSQL(col string, args []string) (string, error) {
	set, err := where.ParseMultiArgs(args)
	if err != nil {
		return "", err
	}
	set, err = MapGeotypes(set)
	if err != nil {
		return "", err
	}
	body := where.WherePart(col, set)

	if body == "" {
		return "", nil
	}

	template := `
AND (
%s
)`
	return fmt.Sprintf(template, body), nil
}

func FixGeotype(token string) (string, error) {
	for _, geotype := range model.GetGeoTypeValues() {
		if strings.EqualFold(token, geotype) {
			return geotype, nil
		}
	}
	return "", fmt.Errorf("%w: %q is not a geotype", ErrInvalidParams, token)
}

// MapGeotypes changes case-insensitive geotypes given in query strings to the specific
// case-sensitive geotypes used in the db.
// For example, "lsoa" will be changed to "LSOA".
// Returns error if a geotype doesn't match at all.
func MapGeotypes(set *where.ValueSet) (*where.ValueSet, error) {
	callback := func(single, low, high *string) (*string, *string, *string, error) {
		var err error

		if single == nil {
			err = fmt.Errorf("%w: cannot have ranges of geotypes", ErrInvalidParams)
		} else {
			var geotype string
			geotype, err = FixGeotype(*single)
			if err == nil {
				single = &geotype
			}
		}
		return single, low, high, err
	}

	return set.Walk(callback)
}
