package where

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// WherePart returns the part of the where clause between parens, as in:
// (
//	   geography_code IN ( ... )
//     OR
//     geography_code BETWEEN (low, high)
//	   ...
// )
//
// col is the name of the column we are matching (eg, "geography_code" or "category_code").
// set is a ValueSet which contains the single value and ranges returned by ParseMultiArgs.
//
// If set has no single values or ranges, an empty string will be returned.
//
func WherePart(col string, set *ValueSet) string {
	var conditions []string

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

	if len(conditions) == 0 {
		return ""
	}
	return strings.Join(conditions, "    OR\n")
}

// ParseMultiArgs generates a ValueSet from rows= and col= multi value arguments.
//
// XXX return proper error types once errors are factored out of geodata
func ParseMultiArgs(args []string) (*ValueSet, error) {
	set := NewValueSet()
	for _, arg := range args {

		// each argument may have many values or ranges separated by commas
		tokens := strings.Split(arg, ",")
		for _, token := range tokens {

			// each token may be a single value or a range
			if !strings.Contains(token, "...") {
				if token == "" {
					return nil, errors.New("value must not be empty")
				}
				set.AddSingle(token)

			} else {
				r := strings.Split(token, "...")
				if len(r) != 2 {
					return nil, errors.New("range must be low...high")
				}
				if r[0] == "" || r[1] == "" {
					return nil, errors.New("range values must not be empty")
				}
				set.AddRange(r[0], r[1])
			}
		}
	}
	return set, nil
}
