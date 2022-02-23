package where

import (
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
