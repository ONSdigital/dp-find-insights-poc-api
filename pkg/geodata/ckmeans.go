package geodata

import (
	"context"
	"fmt"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/timer"
	"github.com/jtrim-ons/ckmeans/pkg/ckmeans"
	"github.com/lib/pq"
)

func getCKmeansSQL(geotypes []string, cats []string) string {
	// sanitise inputs
	var safeGeotypes []string
	for _, geotype := range geotypes {
		safeGeotypes = append(safeGeotypes, pq.QuoteLiteral(geotype))
	}
	var safeCats []string
	for _, cat := range cats {
		safeCats = append(safeCats, pq.QuoteLiteral(cat))
	}
	// return sql
	return fmt.Sprintf(`
SELECT
	geo_type.name
    , geo_metric.metric
	, nomis_category.long_nomis_code
	, geo.id
FROM
    geo,
    geo_type,
    nomis_category,
    geo_metric,
    data_ver

-- the geo_type we are interested in
WHERE geo_type.name in (%s) 

-- all geocodes in this type
AND geo.type_id = geo_type.id
AND geo.valid

-- the category we are interested in
AND nomis_category.long_nomis_code IN (%s)
AND nomis_category.year = $1

-- metrics for these geocodes and category
AND geo_metric.geo_id = geo.id
AND geo_metric.category_id = nomis_category.id

-- only pick metrics for census year / version2.2
AND data_ver.id = geo_metric.data_ver_id
AND data_ver.census_year = nomis_category.year
AND data_ver.ver_string = '2.2'
`,
		strings.Join(safeGeotypes, ", "),
		strings.Join(safeCats, ", "),
	)
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

func (app *Geodata) CKmeans(ctx context.Context, year int, cat string, geotype string, k int, divide_by string) (map[string]map[string][]float64, error) {
	// assemble (hopefully) sanitised sql
	geotypes := strings.Split(geotype, ",")
	categories := strings.Split(cat, ",")
	if divide_by != "" {
		categories = append(categories, divide_by)
	}
	sql := getCKmeansSQL(geotypes, categories)

	// query for data
	t := timer.New("query")
	t.Start()
	rows, err := app.db.DB().QueryContext(
		ctx,
		sql,
		year,
	)
	if err != nil {
		return nil, err
	}
	t.Stop()
	t.Print()
	defer rows.Close()

	// initialise container
	metrics := make(map[string]map[string]map[int]float64)
	for _, category := range categories {
		metrics[category] = make(map[string]map[int]float64)
		for _, geotype := range geotypes {
			metrics[category][geotype] = make(map[int]float64)
		}
	}

	// scan data from rows
	tnext := timer.New("next")
	tscan := timer.New("scan")
	var nmetrics int
	for {
		tnext.Start()
		ok := rows.Next()
		tnext.Stop()
		if !ok {
			if nmetrics == 0 {
				return nil, fmt.Errorf("No data found for %s: %w", strings.Join(categories, ", "), ErrNoContent)
			}
			break
		}
		var (
			geotype string
			metric  float64
			cat     string
			geoID   int
		)
		tscan.Start()
		err := rows.Scan(&geotype, &metric, &cat, &geoID)
		tscan.Stop()
		if err != nil {
			return nil, err
		}
		metrics[cat][geotype][geoID] = metric
		nmetrics++
	}
	tnext.Print()
	tscan.Print()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// prepare container for ckmeans input
	dataForBreaks := make(map[string]map[string][]float64)
	for _, category := range categories {
		if category != divide_by {
			dataForBreaks[category] = make(map[string][]float64)

		}
	}
	// we're nesting geotype under categories, but more efficient to go through by geotype
	for _, geotype := range geotypes {
		if divide_by != "" {
			for geoID, metricDenominator := range metrics[divide_by][geotype] {
				for _, numerator := range categories {
					// NB skip the denominator
					if numerator != divide_by {
						metricNumerator, prs := metrics[numerator][geotype][geoID]
						if !prs {
							return nil, fmt.Errorf("Incomplete data for category %s: %w", numerator, ErrPartialContent)
						}
						dataForBreaks[numerator][geotype] = append(dataForBreaks[numerator][geotype], metricNumerator/metricDenominator)
					}
				}
			}
		} else {
			for cat, metricsCat := range metrics {
				for _, metricCat := range metricsCat[geotype] {
					dataForBreaks[cat][geotype] = append(dataForBreaks[cat][geotype], metricCat)
				}
			}
		}
	}
	// get breaks
	breaks := make(map[string]map[string][]float64)
	for _, cat := range categories {
		// NB skip the denominator
		if cat != divide_by {
			breaks[cat] = make(map[string][]float64)
			for _, geotype := range geotypes {
				catBreaks, err := getBreaks(dataForBreaks[cat][geotype], k)
				if err != nil {
					return nil, err
				}
				breaks[cat][geotype] = catBreaks
			}
		}
	}
	return breaks, nil
}
