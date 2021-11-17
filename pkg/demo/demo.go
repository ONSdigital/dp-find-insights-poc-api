package demo

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/timer"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/where"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

type Demo struct {
	db *database.Database
}

func New(db *database.Database) (*Demo, error) {
	return &Demo{
		db: db,
	}, nil
}

func (app *Demo) Query(ctx context.Context, dataset string, rows, cols []string) (string, error) {
	if dataset == "skinny" {
		return app.skinnyQuery(ctx, rows, cols)
	}
	cols = gatherTokens(cols)
	return app.tableQuery(ctx, dataset, rows, cols)
}

func (app *Demo) skinnyQuery(ctx context.Context, rows, cols []string) (string, error) {
	clause, err := where.SkinnyWhere(rows, cols)
	if err != nil {
		return "", err
	}
	return clause, nil
}

// query runs a SQL query against the db and returns the resulting CSV as a string.
// dataset is the name of the table to query.
// rowspec is not used at present.
// colspec is a list of columns to include in the result. Empty means all columns.
//
func (app *Demo) tableQuery(ctx context.Context, dataset string, rowspec, colspec []string) (string, error) {

	// check allow-list for valid table

	if !validTable(dataset) {
		log.Println("invalid table: " + dataset)
		return "", errors.New("invalid table")
	}

	// parse all the row= variables
	rowvalues, err := where.ParseRows(rowspec)
	if err != nil {
		return "", err
	}

	// We use a string as the output buffer for now.
	// Might hit size limits, so investigate if there can be some kind of streaming output.
	var body strings.Builder
	body.Grow(1000000)

	// Set up to write CSV rows to output buffer.
	w := csv.NewWriter(&body)

	// Construct SQL query.
	// XXX must escape or quote identifiers in the sql statement below XXX
	// Identifiers must be quoted/escaped, but pq.QuoteLiteral and pq.QuoteIdentifier do not work.
	// They surround the resulting strings with " and ', respectively.
	// And quoted strings do not work for some reason:
	// 	"problem with query: pq: relation \"atlas2011.qs101ew\" does not exist"
	//
	var colstring string
	if len(colspec) == 0 {
		colstring = "*"
	} else {
		colstring = strings.Join(colspec, ",")
	}
	sql := fmt.Sprintf(
		`SELECT %s FROM %s %s`,
		colstring,
		dataset,
		where.Clause(rowvalues),
	)

	// log SQL
	fmt.Printf("sql: %s\n", sql)

	// Query the db.
	t := timer.New("query")
	t.Start()
	rows, err := app.db.DB().QueryContext(ctx, sql)
	if err != nil {
		return "", err
	}
	t.Stop()
	t.Print()
	defer rows.Close()

	// Print column names as first row of CSV.
	names, err := rows.Columns()
	if err != nil {
		return "", err
	}
	w.Write(names)

	// Create slice of strings to hold each column value.
	// The Scan mechanism is awkward when we don't know what columns we are
	// getting back.
	// So we make a slice of interface{}'s as expected by Scan, but make sure
	// their concrete types are pointers to strings.
	// This causes the pq library to cast all columns to strings, which is good
	// enough for us now.
	values := make([]interface{}, len(names))
	for i := range values {
		var s string
		values[i] = &s
	}

	// Retrieve each row and print it as a CSV row.
	// For this to work, csv Write has to be given a []string.
	cols := make([]string, 0, len(names))
	tnext := timer.New("next")
	tscan := timer.New("scan")
	twrite := timer.New("write")
	for {
		tnext.Start()
		ok := rows.Next()
		tnext.Stop()
		if !ok {
			break
		}

		tscan.Start()
		err := rows.Scan(values...)
		tscan.Stop()
		if err != nil {
			return "", err
		}
		cols = cols[:0]
		for _, value := range values {
			s, ok := value.(*string)
			if !ok {
				cols = append(cols, "<not a string>")
			} else {
				cols = append(cols, *s)
			}
		}
		twrite.Start()
		w.Write(cols)
		twrite.Stop()
	}
	tnext.Print()
	tscan.Print()
	twrite.Print()

	// check if we stopped because of an error
	if err := rows.Err(); err != nil {
		return "", err
	}

	w.Flush()
	return body.String(), err
}

// gatherTokens collects and combines multiple query parameters.
// For example, turns rows=a,b&rows=c into [a,b,c]
//
func gatherTokens(values []string) []string {
	var tokens []string
	for _, value := range values {
		t := strings.Split(value, ",")
		for _, s := range t {
			if s != "" {
				tokens = append(tokens, s)
			}
		}
	}
	return tokens
}

// XXX probably needs a better solution

func validTable(dataset string) bool {

	m := map[string]bool{
		"atlas2011.qs101ew": true,
		"atlas2011.qs103ew": true,
		"atlas2011.qs104ew": true,
		"atlas2011.qs105ew": true,
		"atlas2011.qs106ew": true,
		"atlas2011.qs108ew": true,
		"atlas2011.qs110ew": true,
		"atlas2011.qs111ew": true,
		"atlas2011.qs112ew": true,
		"atlas2011.qs113ew": true,
		"atlas2011.qs114ew": true,
		"atlas2011.qs115ew": true,
		"atlas2011.qs116ew": true,
		"atlas2011.qs117ew": true,
		"atlas2011.qs118ew": true,
		"atlas2011.qs119ew": true,
		"atlas2011.qs201ew": true,
		"atlas2011.qs202ew": true,
		"atlas2011.qs203ew": true,
		"atlas2011.qs204ew": true,
		"atlas2011.qs205ew": true,
		"atlas2011.qs208ew": true,
		"atlas2011.qs210ew": true,
		"atlas2011.qs211ew": true,
		"atlas2011.qs301ew": true,
		"atlas2011.qs302ew": true,
		"atlas2011.qs303ew": true,
		"atlas2011.qs401ew": true,
		"atlas2011.qs402ew": true,
		"atlas2011.qs403ew": true,
		"atlas2011.qs404ew": true,
		"atlas2011.qs405ew": true,
		"atlas2011.qs406ew": true,
		"atlas2011.qs407ew": true,
		"atlas2011.qs408ew": true,
		"atlas2011.qs409ew": true,
		"atlas2011.qs410ew": true,
		"atlas2011.qs411ew": true,
		"atlas2011.qs412ew": true,
		"atlas2011.qs413ew": true,
		"atlas2011.qs414ew": true,
		"atlas2011.qs415ew": true,
		"atlas2011.qs416ew": true,
		"atlas2011.qs417ew": true,
		"atlas2011.qs418ew": true,
		"atlas2011.qs419ew": true,
		"atlas2011.qs420ew": true,
		"atlas2011.qs421ew": true,
		"atlas2011.qs501ew": true,
		"atlas2011.qs502ew": true,
		"atlas2011.qs601ew": true,
		"atlas2011.qs602ew": true,
		"atlas2011.qs603ew": true,
		"atlas2011.qs604ew": true,
		"atlas2011.qs605ew": true,
		"atlas2011.qs606ew": true,
		"atlas2011.qs607ew": true,
		"atlas2011.qs608ew": true,
		"atlas2011.qs609ew": true,
		"atlas2011.qs610ew": true,
		"atlas2011.qs611ew": true,
		"atlas2011.qs612ew": true,
		"atlas2011.qs613ew": true,
		"atlas2011.qs701ew": true,
		"atlas2011.qs702ew": true,
		"atlas2011.qs703ew": true,
		"atlas2011.qs801ew": true,
		"atlas2011.qs802ew": true,
		"atlas2011.qs803ew": true,
	}

	return m[dataset]

}
