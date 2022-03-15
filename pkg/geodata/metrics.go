package geodata

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/cantabular"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/table"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/timer"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/where"
	"github.com/ONSdigital/dp-find-insights-poc-api/sentinel"
	"github.com/lib/pq"
)

// Retrieve metrics from postgres.
func (app *Geodata) PGMetrics(ctx context.Context, year int, geocodes []string, catset *where.ValueSet, include []string, censustable string) ([]byte, error) {
	sql, include, err := app.metricsSQL(ctx, year, geocodes, catset, include, censustable)
	if err != nil {
		return nil, err
	}

	tbl := table.New()

	var body bytes.Buffer
	body.Grow(1000000)

	t := timer.New("query")
	t.Start()
	rows, err := app.db.DB().QueryContext(ctx, sql)
	if err != nil {
		return nil, err
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
				return nil, fmt.Errorf("%w: limit is %d", sentinel.ErrTooManyMetrics, app.maxMetrics)
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
			return nil, err
		}

		tbl.SetCell(geo, geotype, cat, value)
	}
	tnext.Print()
	tscan.Print()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if nmetrics == 0 {
		return nil, sentinel.ErrNoContent
	}

	tgen := timer.New("generate")
	tgen.Start()
	err = tbl.Generate(&body, include)
	tgen.Stop()
	tgen.Print()
	if err != nil {
		return nil, err
	}

	return body.Bytes(), nil
}

// XXX mv geoCondition into a generic function in where package, or parse geocodes as a valueset
func (app *Geodata) metricsSQL(ctx context.Context, year int, geocodes []string, catset *where.ValueSet, include []string, censustable string) (string, []string, error) {
	// construct AND geo.code IN (...)
	geoCondition := fmt.Sprintf(
		"AND geo.code IN (%s)",
		quoteCodes(geocodes),
	)

	// construct WHERE condition for categories
	catConditions, err := categorySQL(catset, censustable)
	if err != nil {
		return "", nil, err
	}

	// construct additional conditions for censustable / short_nomis_code
	censustableFromSQL, censustableAndSQL := censusTableFromAndSQL(censustable)

	// construct SQL
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
	-- geo conditions
%s
	-- censustable conditions
%s
AND geo_metric.geo_id = geo.id
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = %d
AND data_ver.ver_string = '2.2'
AND nomis_category.id = geo_metric.category_id
AND nomis_category.year = data_ver.census_year
	-- category conditions;
%s
`

	sql := fmt.Sprintf(
		template,
		censustableFromSQL,
		geoCondition,
		censustableAndSQL,
		year,
		catConditions,
	)

	return sql, include, nil
}

func quoteCodes(geocodes []string) string {
	var quoted []string
	for _, code := range geocodes {
		quoted = append(quoted, pq.QuoteLiteral(code))
	}
	return strings.Join(quoted, ",")
}

// Retrieve metrics from Cantabular.
func (app *Geodata) CantabularMetrics(ctx context.Context, geocodes []string, catset *where.ValueSet, geotype string) ([]byte, error) {
	if app.cant == nil {
		return nil, fmt.Errorf("%w: cantabular not enabled", sentinel.ErrNotSupported)
	}

	// the current cantabular queries accept a single category code
	if len(catset.Singles) != 1 {
		return nil, fmt.Errorf("%w: cantabular queries only accept a single category code", sentinel.ErrInvalidParams)
	}
	if len(catset.Ranges) != 0 {
		return nil, fmt.Errorf("%w: cantabular queries do not accept category code ranges", sentinel.ErrInvalidParams)
	}
	if geotype == "" {
		return nil, fmt.Errorf("%w: cantabular queries require geotype", sentinel.ErrMissingParams)
	}

	geoq, catsQL, values, err := app.cant.QueryMetricFilter(ctx, "", strings.Join(geocodes, ","), geotype, catset.Singles[0])
	if err != nil {
		return nil, err
	}
	return []byte(cantabular.ParseMetric(geoq, catsQL, values)), nil
}
