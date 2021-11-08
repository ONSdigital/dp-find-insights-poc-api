package where

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// By default all rows are returned.
// This can be limited by providing one or more rows= query variables.
//
// Each rows variable looks like this:
//
//	rows=col:rowspec[,rowspec]
//
// Where col is a column name, and rowspec is one of:
//
//	single column value, like E01000001
//
// or a range, like E01000003...E01000006
//
// For example:
//
//	rows=geography_code:E01000001,E01000003...E01000006
//
// Multiple rows= variables can be provided for the same column.
// In this case, it is as if all the specs were concatenated.
//
// In the underlying SQL, each single column value is translated to
// an IN condition, and each range is translated to a BETWEEN condition,
// all combined with an OR.
//
// For the example above:
//
//	WHERE geography_code IN ( 'E01000001' )
//	OR geography_code BETWEEN 'E01000003' AND 'E01000006'
//

// Rows holds the values and ranges for each named column.
// The column name is the map key.
type Rows map[string]*ValueSet

// ValueSet collects all the single values and ranges for a column.
type ValueSet struct {
	Singles []string      // list of single values; becomes IN
	Ranges  []*ValueRange // list of value ranges; becomes BETWEEN
}

// ValueRange holds the low and high values for a range.
type ValueRange struct {
	Low  string
	High string
}

// ParseRows takes all raw rows= values and constructs the Rows map
// of values and ranges for each named column.
//
func ParseRows(rows []string) (Rows, error) {
	result := make(Rows)

	for _, q := range rows {

		// split into col and specs on colon
		parts := strings.Split(q, ":")
		if len(parts) != 2 {
			return nil, errors.New("expected col:rowspec[,...]")
		}
		col := parts[0]

		// get ValueSet for this column
		set, ok := result[col]
		if !ok {
			set = &ValueSet{}
			result[col] = set
		}

		// split specs on comma
		rowspecs := strings.Split(parts[1], ",")

		// determine if each spec is a single or a range
		// and append to the set's singles or ranges list.
		for _, rowspec := range rowspecs {
			if len(rowspec) == 0 {
				return nil, errors.New("rowspec must not be empty")
			}
			if !strings.Contains(rowspec, "...") {
				set.Singles = append(set.Singles, rowspec)
			} else {
				parts := strings.Split(rowspec, "...")
				if len(parts) != 2 {
					return nil, errors.New("range must be low...high")
				}
				if parts[0] == "" || parts[1] == "" {
					return nil, errors.New("range values must not be empty")
				}
				vrange := &ValueRange{
					Low:  parts[0],
					High: parts[1],
				}
				set.Ranges = append(set.Ranges, vrange)
			}
		}
	}

	return result, nil
}

func Clause(rows Rows) string {
	var conditions []string

	for col, set := range rows {
		// construct this column's IN condition from Singles
		if len(set.Singles) > 0 {
			var qvalues []string // sql-quoted values
			for _, single := range set.Singles {
				qvalues = append(qvalues, pq.QuoteLiteral(single))
			}
			condition := fmt.Sprintf(
				"%s IN ( %s )",
				col,
				strings.Join(qvalues, ","),
			)
			conditions = append(conditions, condition)
		}

		// construct this columns's BETWEENs from Ranges
		for _, vrange := range set.Ranges {
			condition := fmt.Sprintf(
				"%s BETWEEN %s AND %s",
				col,
				pq.QuoteLiteral(vrange.Low),
				pq.QuoteLiteral(vrange.High),
			)
			conditions = append(conditions, condition)
		}
	}

	// there may not be any conditions if rows was empty
	if len(conditions) == 0 {
		return ""
	}

	// join conditions with OR and stick a WHERE in front
	return "WHERE " + strings.Join(conditions, " OR ")
}
