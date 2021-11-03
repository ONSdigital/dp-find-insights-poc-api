package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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
		`SELECT %s FROM %s`,
		colstring,
		dataset,
	)

	// Query the db.
	rows, err := app.db.QueryContext(ctx, sql)
	if err != nil {
		return "", err
	}
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

	// check if we stopped because of an error
	if err := rows.Err(); err != nil {
		return "", err
	}

	w.Flush()
	return body.String(), err
}

func (app *App) Handler(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
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

func main() {
	app := NewApp()
	lambda.Start(app.Handler)
}
