package table

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/sentinel"
	"github.com/lib/pq"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
)

// A Table holds column names and types used when generating SQL statements.
type Table struct {
	name     string                  // name of table
	colnames []string                // sorted column names, not including wkb_geometry
	coltypes map[string]reflect.Kind // "kind" of each column, indexed by column name
	keepcase bool                    // true if column names should preserve case as found in features
	writer   io.Writer               // print SQL to this writer
}

// New creates a Table that holds column names and types extracted from features.
// Generated SQL will be written to w.
func New(name string, features []geojson.Feature, keepcase bool, w io.Writer) (*Table, error) {
	if name == "" {
		return nil, sentinel.ErrTableName
	}

	// extract column names and types from features
	coltypes, err := getColumnTypes(features)
	if err != nil {
		return nil, err
	}

	// keep sorted list of column names
	var names []string
	for name := range coltypes {
		names = append(names, name)
	}
	sort.Strings(names)

	return &Table{
		name:     name,
		colnames: names,
		coltypes: coltypes,
		keepcase: keepcase,
		writer:   w,
	}, nil
}

// CreateTable produces SQL CREATE TABLE with columns for each property in features,
// and a column for wkb_geometry.
func (tab *Table) CreateTable() error {
	var fragments []string
	for _, name := range tab.colnames {
		var sqltype string
		switch tab.coltypes[name] {
		case reflect.String:
			sqltype = "character varying"
		case reflect.Float64:
			sqltype = "double precision"
		default:
			return fmt.Errorf("%w: %s kind must be string or float64", sentinel.ErrUnusableType, sqltype)
		}
		if !tab.keepcase {
			name = strings.ToLower(name)
		}
		fragments = append(fragments, fmt.Sprintf("%s %s", pq.QuoteIdentifier(name), sqltype))
	}

	// The column name wkb_geometry name is for compatibility with previous ogr2ogr usage.
	const template = `
CREATE EXTENSION IF NOT EXISTS postgis;
DROP TABLE IF EXISTS %s;
CREATE TABLE %s (
    %s,
    wkb_geometry geometry(Geometry,4326)
);
`

	_, err := fmt.Fprintf(
		tab.writer,
		template,
		tab.quotedTableName(),
		tab.quotedTableName(),
		strings.Join(fragments, ",\n    "),
	)
	return err
}

// InsertRows produces a SQL INSERT statements for each feature.
func (tab *Table) InsertRows(features []geojson.Feature) error {
	const template = `
INSERT INTO %s (
    %s,
    wkb_geometry
) VALUES (
    %s,
    ST_GeomFromText('%s',4326)
);
`

	// create SQL fragment holding column names
	namesFrag := strings.Join(tab.quotedColNames(), ",\n    ")

	// produce an INSERT for each feature
	for _, feature := range features {
		// generate SQL-safe string representations of each property value
		sqlValues, err := tab.genSQLvalues(feature.Properties)
		if err != nil {
			return err
		}

		// order values by their column name
		var orderedValues []string
		for _, name := range tab.colnames {
			orderedValues = append(orderedValues, sqlValues[name])
		}

		// join ordered values into SQL fragment
		valuesFrag := strings.Join(orderedValues, ",\n    ")

		// generate WKT string from feature geometry
		wktstring, err := wkt.Marshal(feature.Geometry)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(
			tab.writer,
			template,
			tab.quotedTableName(),
			namesFrag,
			valuesFrag,
			wktstring,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyRows produces SQL to insert rows for each feature using COPY FROM.
func (tab *Table) CopyRows(features []geojson.Feature) error {
	const template = `
COPY %s (
    %s,
    wkb_geometry
)
FROM STDIN;
`

	// create SQL fragment holding column names
	namesFrag := strings.Join(tab.quotedColNames(), ",\n    ")

	// produce COPY FROM
	_, err := fmt.Fprintf(
		tab.writer,
		template,
		tab.quotedTableName(),
		namesFrag,
	)
	if err != nil {
		return err
	}

	// produce line for each feature
	for _, feature := range features {
		// generate COPY-safe string representations of each property value
		copyValues, err := tab.genCOPYvalues(feature.Properties)
		if err != nil {
			return err
		}

		// build list of COPY-safe values ordered by column name
		var orderedValues []string
		for _, name := range tab.colnames {
			orderedValues = append(orderedValues, copyValues[name])
		}

		// generate WKT string from feature geometry
		wktstring, err := wkt.Marshal(feature.Geometry)
		if err != nil {
			return err
		}

		// produce COPY FROM line for this feature
		_, err = fmt.Fprintf(
			tab.writer,
			"%s\t%s\n",
			strings.Join(orderedValues, "\t"),
			wktstring,
		)
		if err != nil {
			return err
		}

	}

	// and final line to terminate the data section
	_, err = fmt.Fprintf(tab.writer, "\\.\n")
	return err
}

// getColumnTypes scans the properties of every feature to determine column names and types
// that will be needed when creating the table.
// It is an error for columns to have different types in different features.
func getColumnTypes(features []geojson.Feature) (map[string]reflect.Kind, error) {
	cols := map[string]reflect.Kind{}

	for _, feature := range features {
		for name, value := range feature.Properties {
			thisKind := reflect.ValueOf(value).Kind()
			kind, seen := cols[name]
			if !seen {
				cols[name] = thisKind
			} else if thisKind != kind {
				return nil, fmt.Errorf("%w: property %s: %s and %s", sentinel.ErrInconsistentTypes, name, kind, thisKind)
			}
		}
	}
	return cols, nil
}

// quotedTableName returns the table name in SQL-safe form.
func (tab *Table) quotedTableName() string {
	return pq.QuoteIdentifier(tab.name)
}

// quotedColNames returns a slice of column names, each SQL-identifier quoted.
func (tab *Table) quotedColNames() []string {
	var result []string
	for _, name := range tab.colnames {
		if !tab.keepcase {
			name = strings.ToLower(name)
		}
		result = append(result, pq.QuoteIdentifier(name))
	}
	return result
}

// genSQLvalues creates a map of SQL-safe column values as strings.
// The map is keyed on column names.
func (tab *Table) genSQLvalues(properties map[string]interface{}) (map[string]string, error) {
	values := map[string]string{}
	for _, name := range tab.colnames {
		var v interface{}
		var sql string
		v, ok := properties[name]
		if !ok {
			sql = "NULL"
		} else {
			switch v := v.(type) {
			case nil:
				sql = "NULL"
			case string:
				sql = pq.QuoteLiteral(v)
			case float64:
				sql = strconv.FormatFloat(v, 'g', -1, 64)
			default:
				return nil, fmt.Errorf("%w: %s type must be nil, string or float64", sentinel.ErrUnusableType, name)
			}
		}
		values[name] = sql
	}
	return values, nil
}

// genCopyValues creates a map of COPY-safe column values as strings.
// The map is keyed on column names.
func (tab *Table) genCOPYvalues(properties map[string]interface{}) (map[string]string, error) {
	values := map[string]string{}
	for _, name := range tab.colnames {
		var v interface{}
		var copy string
		v, ok := properties[name]
		if !ok {
			copy = `\N`
		} else {
			switch v := v.(type) {
			case nil:
				copy = `\N`
			case string:
				copy = copyEscape(v)
			case float64:
				copy = strconv.FormatFloat(v, 'g', -1, 64)
			default:
				return nil, fmt.Errorf("%w: %s type must be nil, string or float64", sentinel.ErrUnusableType, name)
			}
		}
		values[name] = copy
	}
	return values, nil
}

// copyEscape escapes string s for use in a COPY FROM line.
// See https://www.postgresql.org/docs/10/sql-copy.html
// There is probably a better way to do this.
func copyEscape(s string) string {
	var escapes = []struct {
		from string
		to   string
	}{
		{`\`, `\\`}, // this has to be first
		{"\b", `\b`},
		{"\f", `\f`},
		{"\n", `\n`},
		{"\r", `\r`},
		{"\t", `\t`},
		{"\v", `\v`},
	}

	for _, esc := range escapes {
		s = strings.ReplaceAll(s, esc.from, esc.to)
	}
	return s
}
