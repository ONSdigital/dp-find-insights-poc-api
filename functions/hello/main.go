package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/timer"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/where"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	_ "github.com/lib/pq"
)

type App struct {
	aws *session.Session               // AWS session
	sm  *secretsmanager.SecretsManager // secrets manager client
	db  *sql.DB

	// It's hard to return usable errors from a lambda main function.
	// So any errors in main set err here and Handler checks it first.
	errmsg string
	err    error
}

func errorResponse(status int, logmsg string, err error) *events.APIGatewayProxyResponse {
	if err != nil {
		logmsg = logmsg + ": " + err.Error()
	}
	log.Println(logmsg)

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusText(status),
		Message: logmsg,
	}

	body, err := json.MarshalIndent(&response, "", "    ")
	if err != nil {
		log.Println(err.Error())
		body = []byte(err.Error())
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
	}
}

func clientResponse(msg string) *events.APIGatewayProxyResponse {
	return errorResponse(http.StatusBadRequest, msg, nil)
}

func NewApp() *App {
	app := &App{}

	// Set up AWS session
	cfg := aws.NewConfig()
	sess, err := session.NewSession(cfg)
	if err != nil {
		return &App{
			errmsg: "cannot get session",
			err:    err,
		}
	}
	app.aws = sess

	// Create Secrets Manager client
	app.sm = secretsmanager.New(app.aws, aws.NewConfig())

	// Look up postgres password
	secret, err := app.getSecret()
	if err != nil {
		log.Println("getSecret returned error")
		return &App{
			errmsg: "cannot get secret",
			err:    err,
		}
	}

	// totally environment variables for now, so no connection string
	db, err := sql.Open("postgres", fmt.Sprintf("password=%s", secret))
	if err != nil {
		return &App{
			errmsg: "cannot Open postgres",
			err:    err,
		}
	}

	//	ctx := context.Background()
	//	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	//	defer cancel()

	//	// check db visible
	//	err = db.PingContext(ctx)
	//	if err != nil {
	//		return &App{
	//			errmsg: "cannot Ping postgres",
	//			err:    err,
	//		}
	//	}

	app.db = db
	return app
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

// query runs a SQL query against the db and returns the resulting CSV as a string.
// dataset is the name of the table to query.
// rowspec is not used at present.
// colspec is a list of columns to include in the result. Empty means all columns.
//
func (app *App) query(ctx context.Context, dataset string, rowspec, colspec []string) (string, error) {

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
	rows, err := app.db.QueryContext(ctx, sql)
	if err != nil {
		return "", err
	}
	t.Stop()
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
	t = timer.New("scans")
	for rows.Next() {
		err := rows.Scan(values...)
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
		w.Write(cols)
	}
	t.Stop()

	// check if we stopped because of an error
	if err := rows.Err(); err != nil {
		return "", err
	}

	w.Flush()
	return body.String(), err
}

func (app *App) Handler(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	// log req (probably too verbose)
	fmt.Printf("req: %#v\n", req)

	// return any init errors
	if app.err != nil {
		return errorResponse(http.StatusInternalServerError, app.errmsg, app.err), nil
	}

	// grab parameters from path and query string
	//
	dataset := req.PathParameters["dataset"]
	if dataset == "" {
		return clientResponse("missing dataset path parameter"), nil
	}
	// empty list means ALL
	rows := gatherTokens(req.MultiValueQueryStringParameters["rows"])
	cols := gatherTokens(req.MultiValueQueryStringParameters["cols"])

	body, err := app.query(ctx, dataset, rows, cols)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "problem with query", err), nil
	}

	response := &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       body,
	}

	return response, nil
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

func main() {
	app := NewApp()
	lambda.Start(app.Handler)
}
