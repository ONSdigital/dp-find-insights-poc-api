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
	"fmt"
	"io"
	"sort"
)

const (
	ColGeographyCode = "geography_code"
	ColGeotype       = "geotype"
)

type Geocode string // eg "E07000107"
type Geotype string // eg "LSOA", "LAD", ...
type Catcode string // eg "QS412EW0001"

type Table struct {
	geocodes map[Geocode]bool // geography codes seen
	catcodes map[Catcode]bool // category codes seen
	areas    map[Geocode]area // areas indexed by geography code
}

type area struct {
	geotype Geotype
	metrics map[Catcode]float64
}

// New creates a new table.
//
func New() *Table {
	return &Table{
		geocodes: map[Geocode]bool{},
		catcodes: map[Catcode]bool{},
		areas:    map[Geocode]area{},
	}
}

// SetCell sets the value of a cell on the row matching geocode and geotype, and the column matching colname.
// New rows and columns are created dynamically.
// XXX check for duplicate geotype XXX move errors.go to its own pkg/errors
func (tbl *Table) SetCell(geocode, geotype, catcode string, value float64) {
	// Remember the geo and cat codes for when we generate the table
	tbl.geocodes[Geocode(geocode)] = true
	tbl.catcodes[Catcode(catcode)] = true

	// Look up or create area for this geocode
	a, ok := tbl.areas[Geocode(geocode)]
	if !ok {
		a = area{
			geotype: Geotype(geotype),
			metrics: map[Catcode]float64{},
		}
		tbl.areas[Geocode(geocode)] = a
	}

	// Set this geo/type/cat value
	a.metrics[Catcode(catcode)] = value
}

// Generate produces a CSV version of the table on w.
// It doesn't close w.
//
// include is a list of non-category columns to include in the output table.
// Currently supported values are "geography_code" and "geotype".
//
func (tbl *Table) Generate(w io.Writer, include []string) error {
	// sort the geography codes we have seen
	geocodes := sort.StringSlice{}
	for geo := range tbl.geocodes {
		geocodes = append(geocodes, string(geo))
	}
	geocodes.Sort()

	// sort the category codes we have seen
	catcodes := sort.StringSlice{}
	for cat := range tbl.catcodes {
		catcodes = append(catcodes, string(cat))
	}
	catcodes.Sort()

	// note which non-category columns we want to include
	// XXX make it an error on unrecognized columns
	includeGeocode := false
	includeGeotype := false
	for _, col := range include {
		switch col {
		case ColGeographyCode:
			includeGeocode = true
		case ColGeotype:
			includeGeotype = true
		}
	}

	// set up csv output on w
	cw := csv.NewWriter(w)

	// print column headings
	colnames := []string{}
	if includeGeocode {
		colnames = append(colnames, ColGeographyCode)
	}
	if includeGeotype {
		colnames = append(colnames, ColGeotype)
	}
	colnames = append(colnames, catcodes...)
	cw.Write(colnames)

	// pre-allocate slice to hold column values
	row := make([]string, len(colnames)+1)

	// generate table ordered by geocode, with each row ordered by category code
	for _, geocode := range geocodes {
		row = row[:0]
		if includeGeocode {
			row = append(row, geocode)
		}
		if includeGeotype {
			row = append(row, string(tbl.areas[Geocode(geocode)].geotype))
		}

		for _, catcode := range catcodes {
			// Precision may need to be increased if numbers are printed as exponents,
			// or if decimals are rounded
			// See the "specific numeric formatting tests" in table_test.go.
			row = append(row, fmt.Sprintf("%.13g", tbl.areas[Geocode(geocode)].metrics[Catcode(catcode)]))
		}

		cw.Write(row)
	}

	cw.Flush()
	return cw.Error()
}
