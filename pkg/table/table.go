// The table package implements a simplistic 2-dimensional array that can be populated one cell at a time,
// and output as a CSV.
// It is intended to be used to build up a wide table from results of queries on the geo_metric table.
//
// So input that looks like this:
//
//	HERE	COLA	10
//  HERE	COLB	20
//	THERE	COLA	30
//	THERE	COLB	40
//
// Becomes
//
//  GEO		COLA	COLB
//	HERE	10		20
//	THERE	30		40
//
package table

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
)

type Table struct {
	primary string
	geos    map[string]bool               // geography codes seen
	cats    map[string]bool               // category code seen
	rows    map[string]map[string]float64 // indexed by geography code, each row indexed by category code
}

// New creates a new table.
// primary is the column name you want to use for the geography_code column in output.
// primary must not be empty.
//
func New(primary string) (*Table, error) {
	if primary == "" {
		return nil, errors.New("primary must not be blank")
	}

	return &Table{
		primary: primary,
		geos:    map[string]bool{},
		cats:    map[string]bool{},
		rows:    map[string]map[string]float64{},
	}, nil
}

// SetCell sets the value of a cell on the row matching geocode and the column matching colname.
// New rows and columns are created dynamically.
//
func (tbl *Table) SetCell(geocode, catcode string, value float64) {
	// Remember the geo and cat codes for when we generate the table
	tbl.geos[geocode] = true
	tbl.cats[catcode] = true

	// Look up or create row for this geocode
	r, ok := tbl.rows[geocode]
	if !ok {
		r = map[string]float64{}
		tbl.rows[geocode] = r
	}

	// Set this geo/cat value
	r[catcode] = value
}

// Generate produces a CSV version of the table on w.
// It doesn't close w.
//
func (tbl *Table) Generate(w io.Writer) error {
	// sort the geography codes we have seen
	geos := sort.StringSlice{}
	for geo := range tbl.geos {
		geos = append(geos, geo)
	}
	geos.Sort()

	// sort the category codes we have seen
	cats := sort.StringSlice{}
	for cat := range tbl.cats {
		cats = append(cats, cat)
	}
	cats.Sort()

	// set up csv output on w
	cw := csv.NewWriter(w)

	// print column headings
	colnames := append([]string{tbl.primary}, cats...)
	cw.Write(colnames)

	// pre-allocate slice to hold column values
	row := make([]string, len(cats)+1)

	// generate column values in order for each row in order
	for _, geo := range geos {
		row = row[:0]
		row = append(row, geo)

		for _, cat := range cats {
			// Precision may need to be increased if numbers are printed as exponents,
			// or if decimals are rounded
			// See the "specific numeric formatting tests" in table_test.go.
			row = append(row, fmt.Sprintf("%.13g", tbl.rows[geo][cat]))
		}

		cw.Write(row)
	}

	cw.Flush()
	return cw.Error()
}
