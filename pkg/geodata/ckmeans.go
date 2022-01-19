package geodata

import (
	"context"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/timer"
	"github.com/jtrim-ons/ckmeans/pkg/ckmeans"
)

func (app *Geodata) CKmeans(ctx context.Context, cat string, geotype string, k int) ([]float64, error) {
	sql := `
SELECT
    geo_metric.metric
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
AND nomis_category.long_nomis_code = $2
AND nomis_category.year = 2011

-- metrics for these geocodes and category
AND geo_metric.geo_id = geo.id
AND geo_metric.category_id = nomis_category.id

-- only pick metrics for 2011/2.2
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
AND data_ver.ver_string = '2.2'
`

	t := timer.New("query")
	t.Start()
	rows, err := app.db.DB().QueryContext(
		ctx,
		sql,
		geotype,
		cat,
	)
	if err != nil {
		return nil, err
	}
	t.Stop()
	t.Print()
	defer rows.Close()

	tnext := timer.New("next")
	tscan := timer.New("scan")
	var nmetrics int
	var metrics []float64
	for {
		tnext.Start()
		ok := rows.Next()
		tnext.Stop()
		if !ok {
			break
		}

		nmetrics++

		var metric float64
		tscan.Start()
		err := rows.Scan(&metric)
		tscan.Stop()
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}
	tnext.Print()
	tscan.Print()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if nmetrics == 0 {
		return nil, ErrNoContent
	}

	return getBreaks(metrics, k)
}

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

func (app *Geodata) CKmeansRatio(ctx context.Context, cat1 string, cat2 string, geotype string, k int) ([]float64, error) {
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
AND nomis_category.year = 2011

-- metrics for these geocodes and category
AND geo_metric.geo_id = geo.id
AND geo_metric.category_id = nomis_category.id

-- only pick metrics for 2011/2.2
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = 2011
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
		return nil, ErrNoContent
	}
	if nmetricsCat1 != nmetricsCat2 {
		return nil, ErrPartialContent
	}

	var metrics []float64
	for geoID, metricCat1 := range metricsCat1 {
		metricCat2, prs := metricsCat2[geoID]
		if prs == false {
			return nil, ErrPartialContent
		}
		metrics = append(metrics, metricCat1/metricCat2)
	}

	return getBreaks(metrics, k)
}
