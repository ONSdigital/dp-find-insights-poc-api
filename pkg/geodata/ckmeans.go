package geodata

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/timer"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/where"
	"github.com/ONSdigital/dp-find-insights-poc-api/sentinel"
	"github.com/jtrim-ons/ckmeans/pkg/ckmeans"
)

// Chunk holds multiple catgeory data for a single geocode
//
type Chunk struct {
	geocode string
	geotype string
	metrics map[string]float64
}

// CkmeansParser holds data and methods neccessary to parse and process data for ckmeans queries
//
type CkmeansParser struct {
	catcodes   []string
	geotypes   []string
	divideBy   string
	k          int
	rowGeocode string
	rowGeotype string
	nmetrics   int
	metrics    map[string]map[string][]float64
	breaks     map[string]map[string][]float64
	chunk      *Chunk
}

// New creates a new CkmeansParser.
//
func NewCkmeansParser(divideBy string, k int) *CkmeansParser {
	return &CkmeansParser{
		catcodes:   []string{},
		geotypes:   []string{},
		divideBy:   divideBy,
		k:          k,
		rowGeocode: "",
		rowGeotype: "",
		chunk: &Chunk{
			geocode: "",
			geotype: "",
			metrics: map[string]float64{},
		},
		nmetrics: 0,
		metrics:  map[string]map[string][]float64{},
		breaks:   map[string]map[string][]float64{},
	}
}

// ------------------------------------------------- main ----------------------------------------------------------- //

// CKmeans does an 'all rows' query for census data using sql generator from geodata pkg, parses the results by
// category and then by geotype, and gets ckmeans breaks for each geotype in each category. Optionally, if divideBy
// is not blank, will divide all categories by the cateogry indicated by the divideBy value, prior to getting ckmeans
// breaks.
//
func (app *Geodata) CKmeans(ctx context.Context, year int, cat []string, geotype []string, k int, divideBy string) (map[string]map[string][]float64, error) {
	// initialise
	ckparser := NewCkmeansParser(divideBy, k)

	// parse and validate tokens
	if err := ckparser.parseCat(cat); err != nil {
		return nil, err
	}
	if err := ckparser.parseValidateGeotype(geotype); err != nil {
		return nil, err
	}

	// get sql
	sql, err := getCkmeansSQL(ctx, year, ckparser)
	if err != nil {
		return nil, err
	}

	// query for data
	t := timer.New("query")
	t.Start()
	rows, err := app.db.DB().QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	t.Stop()
	t.Print()
	defer rows.Close()

	// scan data from rows
	tnext := timer.New("next")
	tscan := timer.New("scan")
	for {
		tnext.Start()
		ok := rows.Next()
		tnext.Stop()

		// if we've got the end of the rows, check we got any data at all
		if !ok {
			if ckparser.nmetrics == 0 {
				return nil, fmt.Errorf("No data found for %s: %w", strings.Join(ckparser.catcodes, ", "), sentinel.ErrNoContent)
			}
			break
		}

		// consume row data
		if err := ckparser.processRow(rows, tscan); err != nil {
			return nil, err
		}
	}
	tnext.Print()
	tscan.Print()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// process all breaks and return
	if err = ckparser.processBreaks(); err != nil {
		return nil, err
	}
	return ckparser.breaks, nil
}

// ---------------------------------------------- CkmeansParser methods --------------------------------------------- //

// CkmeansParser.parseCat parses single values and combines with split comma-seperated cat values and returns as array.
// Will return error if any cat range values (cat1..cat2) are found (for error handling we need to know
// explicitly beforehand which cats to expect in our results)
//
func (ckparser *CkmeansParser) parseCat(cat []string) error {
	catset, err := where.ParseMultiArgs(cat)
	if err != nil {
		return err
	}
	if catset.Ranges != nil {
		return fmt.Errorf("%w: ckmeans endpoint does not accept range values for cats", sentinel.ErrInvalidParams)
	}
	ckparser.catcodes = catset.Singles
	return nil
}

// CkmeansParser.parseValidateGeotype arses single values and combines with split comma-seperated geotype values
// and returns as array. Any badly-cased geotype values will be corrected, and any unrecognised geotype
// value will cause an error to be returned.
//
func (ckparser *CkmeansParser) parseValidateGeotype(geotype []string) error {
	geoset, err := where.ParseMultiArgs(geotype)
	if err != nil {
		return err
	}
	geoset, err = MapGeotypes(geoset)
	if err != nil {
		return err
	}
	ckparser.geotypes = geoset.Singles
	return nil
}

// CkmeansParser.processRow loads data from a sql.Row into a Chunk containing data from a single geotype and geocode
// (NB this only works because the SQL returned by getCkmeansSQL orders by geocode). If the current row contains data
// from a different geotype or geocode, the Chunk is considered complete and is processed and reset.
//
func (ckparser *CkmeansParser) processRow(rows *sql.Rows, tscan *timer.Timer) error {
	// read data from row
	var rowCatcode string
	var rowMetric float64

	tscan.Start()
	err := rows.Scan(&ckparser.rowGeocode, &ckparser.rowGeotype, &rowCatcode, &rowMetric)
	tscan.Stop()
	if err != nil {
		return err
	}

	// reset chunk with new geo and geotype if this is the first row we've read
	if ckparser.nmetrics == 0 {
		ckparser.resetChunk()
	}

	// if geotype OR geoID changes then we have reached the end of a chunk and should process it.
	// NB - geotype will change WITHOUT geoID changing if we only have one category, so check both
	if ckparser.isChunkComplete() {
		if err := ckparser.processChunk(); err != nil {
			return err
		}
	}

	// otherwise collect results - ordered by geo, so should come in chunks
	ckparser.addToChunk(rowCatcode, rowMetric)
	ckparser.nmetrics++
	return nil
}

// CkmeansParser.processChunk parses data from a Chunk containing required category data for single geotype and geocode
// into the main CkmeansParser.metrics container. If CkmeansParser.divideBy is not blank, all categories will be
// divided by the divideBy category before being stored. Chunk is reset after data has been processed from it.
//
func (ckparser *CkmeansParser) processChunk() error {
	// check divideBy if doing ratios
	metricDenominator, prs := ckparser.chunk.metrics[ckparser.divideBy]
	if !prs && ckparser.divideBy != "" {
		return fmt.Errorf("Incomplete data for category %s: %w", ckparser.divideBy, sentinel.ErrPartialContent)
	}

	// process all other required catcodes
	for _, catcode := range ckparser.catcodes {

		// check catcode has metric
		metricCatcode, prs := ckparser.chunk.metrics[catcode]
		if !prs {
			return fmt.Errorf("Incomplete data for category %s: %w", catcode, sentinel.ErrPartialContent)
		}

		// derive or get metric
		var outputMetric float64
		if ckparser.divideBy != "" {
			// make ratio if doing that
			outputMetric = metricCatcode / metricDenominator
		} else {
			// otherwise just use data as is
			outputMetric = metricCatcode
		}

		// append to metrics
		if _, prs := ckparser.metrics[catcode]; !prs {
			ckparser.metrics[catcode] = map[string][]float64{}
		}
		ckparser.metrics[catcode][ckparser.chunk.geotype] = append(ckparser.metrics[catcode][ckparser.chunk.geotype], outputMetric)
	}

	// reset chunk
	ckparser.resetChunk()
	return nil
}

// CkmeansParser.isChunkComplete returns True if CkmeansParser.Chunk contains data from a different geotype or geocode
// from the sql.Row CkmeansParser is currently processing.
//
func (ckparser *CkmeansParser) isChunkComplete() bool {
	// if geotype OR geoID changes then we have reached the end of a chunk and should process it.
	// NB - geotype will change WITHOUT geoID changing if we only have one category, so check both
	return ckparser.rowGeotype != ckparser.chunk.geotype || ckparser.rowGeocode != ckparser.chunk.geocode
}

// CkmeansParser.addToChunk adds data for a new category code to CkmeansParser.Chunk
//
func (ckparser *CkmeansParser) addToChunk(catcode string, metric float64) {
	ckparser.chunk.metrics[catcode] = metric
}

// CkmeansParser.resetChunk deletes all data from CkmeansParser.Chunk and sets the Chunk geotype and geocode to that of
// the sql.Row that CkmeansParser is currently processing.
//
func (ckparser *CkmeansParser) resetChunk() {
	for k := range ckparser.chunk.metrics {
		delete(ckparser.chunk.metrics, k)
	}
	ckparser.chunk.geocode = ckparser.rowGeocode
	ckparser.chunk.geotype = ckparser.rowGeotype
}

// CkmeansParser.processBreaks crawls through CkmeansParser.metrics and runs getBreaks on each geotype's data for each
// category, storing the results in CkmeansParser.breaks
//
func (ckparser *CkmeansParser) processBreaks() error {
	for _, catcode := range ckparser.catcodes {
		for _, geotype := range ckparser.geotypes {
			catBreaks, err := getBreaks(ckparser.metrics[catcode][geotype], ckparser.k)
			if err != nil {
				return err
			}
			// append to breaks
			if _, prs := ckparser.breaks[catcode]; !prs {
				ckparser.breaks[catcode] = map[string][]float64{}
			}
			ckparser.breaks[catcode][geotype] = catBreaks
		}
	}
	return nil
}

// ---------------------------------------------- functions --------------------------------------------------------- //

// getCkmeansSQL validates supplied arguments and then generates SQL for ckmeans query. This is mostly the same as the
// SQL used for a general rows=all query with a geotype filter, but with an additional clause to order results by
// geocode (this is needed to process data in chunks)
//
func getCkmeansSQL(ctx context.Context, year int, ckparser *CkmeansParser) (string, error) {
	// make 'cols' arg for using CensusQuerySQL (cats plus divide_by, if present)
	cols := make([]string, len(ckparser.catcodes))
	copy(cols, ckparser.catcodes)
	if ckparser.divideBy != "" {
		cols = append(cols, ckparser.divideBy)
	}

	// get sql
	sql, _, err := CensusQuerySQL(
		ctx,
		CensusQuerySQLArgs{
			Year:     year,
			Geos:     []string{"all"},
			Geotypes: ckparser.geotypes,
			Cols:     cols,
		},
	)

	// append ORDER_BY to allow chunking of data
	sql = sql + `
	ORDER BY
		geo.id ASC;
	`
	return sql, err
}

// getBreaks gets k ckmeans clusters from metrics and returns the upper breakpoints for each cluster.
//
func getBreaks(metrics []float64, k int) ([]float64, error) {
	clusters, err := ckmeans.Ckmeans(metrics, k)
	if err != nil {
		return nil, err
	}

	var breaks []float64
	for _, cluster := range clusters {
		bp := cluster[len(cluster)-1]
		breaks = append(breaks, bp)
	}
	return breaks, nil
}

// !!!! DEPRECATED CKMEANSRATIO TO BE REMOVED WHEN FRONT END REMOVES DEPENDENCY ON IT !!!!
//
func (app *Geodata) CKmeansRatio(ctx context.Context, year int, cat1 string, cat2 string, geotype string, k int) ([]float64, error) {
	sql := `
SELECT
    geo_metric.metric
	, nomis_category.long_nomis_code
	, geo.id
FROM
    geo,
    geo_type,
    nomis_category,
    geo_metric,
    data_ver
-- the geo_type we are interested in
WHERE geo_type.name = $1 
-- all geocodes in this type
AND geo.type_id = geo_type.id
AND geo.valid
-- the category we are interested in
AND nomis_category.long_nomis_code IN ($2, $3)
AND nomis_category.year = $4
-- metrics for these geocodes and category
AND geo_metric.geo_id = geo.id
AND geo_metric.category_id = nomis_category.id
-- only pick metrics for census year / version2.2
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = nomis_category.year
AND data_ver.ver_string = '2.2'
`

	t := timer.New("query")
	t.Start()
	rows, err := app.db.DB().QueryContext(
		ctx,
		sql,
		geotype,
		cat1,
		cat2,
		year,
	)

	if err != nil {
		return nil, err
	}
	t.Stop()
	t.Print()
	defer rows.Close()

	tnext := timer.New("next")
	tscan := timer.New("scan")

	var nmetricsCat1 int
	var nmetricsCat2 int
	metricsCat1 := make(map[int]float64)
	metricsCat2 := make(map[int]float64)
	for {
		tnext.Start()
		ok := rows.Next()
		tnext.Stop()
		if !ok {
			break
		}

		var (
			metric float64
			cat    string
			geoID  int
		)
		tscan.Start()
		err := rows.Scan(&metric, &cat, &geoID)
		tscan.Stop()
		if err != nil {
			return nil, err
		}
		if cat == cat1 {
			nmetricsCat1++
			metricsCat1[geoID] = metric
		}
		if cat == cat2 {
			nmetricsCat2++
			metricsCat2[geoID] = metric
		}
	}
	tnext.Print()
	tscan.Print()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if nmetricsCat1 == 0 && nmetricsCat2 == 0 {
		return nil, sentinel.ErrNoContent
	}
	if nmetricsCat1 != nmetricsCat2 {
		return nil, sentinel.ErrPartialContent
	}

	var metrics []float64
	for geoID, metricCat1 := range metricsCat1 {
		metricCat2, prs := metricsCat2[geoID]
		if !prs {
			return nil, sentinel.ErrPartialContent
		}
		metrics = append(metrics, metricCat1/metricCat2)
	}

	return getBreaks(metrics, k)
}
