package where

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// The where clause for queries against the geo_metrics table looks like this:
//
// WHERE (
//	   geography_code IN ( ... )
//     OR
//     geography_code BETWEEN (low, high)
//	   ...
// ) AND (
//	   category_code IN ( ... )
//     OR
//     category_code BETWEEN (low, high)
//     ...
// )
//
// Because of the size of the geo_metrics table, we require at least one geography_code
// and at least one category_code to be specified.
//
//

// SkinnyWhere generates the WHERE clause for queries against the geo_metrics table.
// rows and cols are the multivalue rows= and cols= arguments.
//
func SkinnyWhere(rows, cols []string) (string, error) {
	if len(rows) == 0 || len(cols) == 0 {
		return "", errors.New("rows and cols required")
	}

	geo, err := WherePart("geography_code", rows)
	if err != nil {
		return "", err
	}

	cat, err := WherePart("category_code", cols)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("WHERE (\n%s) AND (\n%s)\n", geo, cat), nil
}

// WherePart returns the part of the where clause (see above) between parens.
// col is the name of the column we are matching (eg, "geography_code" or "category_code").
// args is the list of query string values taken from MultiValueQueryStringParameters.
//
func WherePart(col string, args []string) (string, error) {
	var conditions []string

	set, err := GetValues(args)
	if err != nil {
		return "", err
	}

	if len(set.Singles) > 0 {
		var values []string
		for _, single := range set.Singles {
			values = append(values, pq.QuoteLiteral(single))
		}
		condition := fmt.Sprintf(
			"    %s IN ( %s )\n",
			col,
			strings.Join(values, ", "),
		)
		conditions = append(conditions, condition)
	}

	for _, vrange := range set.Ranges {
		condition := fmt.Sprintf(
			"    %s BETWEEN %s AND %s\n",
			col,
			pq.QuoteLiteral(vrange.Low),
			pq.QuoteLiteral(vrange.High),
		)
		conditions = append(conditions, condition)
	}

	return strings.Join(conditions, "    OR\n"), nil
}

// GetValues generates a ValueSet from rows= and col= multi value arguments.
//
func GetValues(args []string) (*ValueSet, error) {
	set := &ValueSet{}
	for _, arg := range args {

		// each argument may have many values or ranges separated by commas
		tokens := strings.Split(arg, ",")
		for _, token := range tokens {

			// each token may be a single value or a range
			if !strings.Contains(token, "...") {
				if token == "" {
					return nil, errors.New("value must not be empty")
				}
				set.Singles = append(set.Singles, token)

			} else {
				r := strings.Split(token, "...")
				if len(r) != 2 {
					return nil, errors.New("range must be low...high")
				}
				if r[0] == "" || r[1] == "" {
					return nil, errors.New("range values must not be empty")
				}
				vr := &ValueRange{
					Low:  r[0],
					High: r[1],
				}
				set.Ranges = append(set.Ranges, vr)
			}
		}
	}
	return set, nil
}
