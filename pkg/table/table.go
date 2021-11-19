// The table package implements a simplistic 2-dimensional array that can be populated one cell at a time,
// and output as a CSV.
// It is intended to be used to build up a wide table from results of queries on the new "skinny" table.
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
// Expected column names must be known in advance (which they are because they are used in the SQL query),
// and you have to provide some name for the geography_code column (GEO in the example above).
//
// This implementation can be optimised when we know more about the problem:
//	* we hold all row and column until the very end; if we know geography codes will be contiguous on input,
//    we could conceivably only hold a row at a time and print each row as it is finished
//  * the simplistic [][]string data structure requires linear searches through row and column keys to find a match;
//    with big datasets this might be inefficient
//
package table

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

type Table struct {
	colnames []string
	rows     [][]string
}

// New creates a new table.
// primary is the column name you want to use for the geography_code column in output.
// colnames is the list of columns (categories or attributes) you expect to see on input.
// It is an error to have an empty primary or colname, and it is an error to have duplicate
// column names.
//
func New(primary string, colnames []string) (*Table, error) {
	if primary == "" {
		return nil, errors.New("primary must not be blank")
	}
	if len(colnames) == 0 {
		return nil, errors.New("table requires at least one column")
	}

	fmt.Fprintf(os.Stderr, "table.New colnames=%+v\n", colnames)
	seen := map[string]bool{}
	for _, name := range colnames {
		if name == "" {
			return nil, errors.New("blank column name")
		}
		if seen[name] {
			return nil, errors.New("duplicate column names")
		}
		seen[name] = true
	}

	tbl := &Table{
		colnames: append([]string{primary}, colnames...),
	}

	fmt.Fprintf(os.Stderr, "tbl=%+v\n", tbl)
	return tbl, nil
}

// SetCell sets the value of a cell on the row matching geocode and the column matching colname.
// New rows are created dynamically, but colname must be one specified in the call to New.
//
// value is a string because we are building up a data structure that can be used directly by the csv library.
//
func (tbl *Table) SetCell(geocode, colname, value string) error {
	i := tbl.colIndex(colname)
	if i == -1 {
		return fmt.Errorf("invalid column name %q", colname)
	}
	tbl.findRow(geocode)[i] = value
	return nil
}

// Generate produces a CSV version of the table on w.
// It doesn't close w.
//
func (tbl *Table) Generate(w io.Writer) error {
	cw := csv.NewWriter(w)

	cw.Write(tbl.colnames)
	cw.WriteAll(tbl.rows)
	cw.Flush()
	return cw.Error()
}

// find Row locates the row corresponding to geocode.
// If the row is not found, a new one is created.
// The []string that is returned corresponds to the row in the table, with element 0 set to geocode and all the other elements empty strings.
//
func (tbl *Table) findRow(geocode string) []string {
	for _, row := range tbl.rows {
		if row[0] == geocode {
			return row
		}
	}
	newrow := make([]string, len(tbl.colnames))
	newrow[0] = geocode
	tbl.rows = append(tbl.rows, newrow)
	return newrow
}

// colIndex returns the index of the colunm named name, or -1 if the column doesn't exist.
//
func (tbl *Table) colIndex(name string) int {
	for i, s := range tbl.colnames {
		if s == name {
			return i
		}
	}
	return -1
}
