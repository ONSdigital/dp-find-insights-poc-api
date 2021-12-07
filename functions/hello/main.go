package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/aws"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
)

type App struct {
	d *geodata.Geodata

	// It's hard to return usable errors from a lambda main function.
	// So any errors in main set err here and Handler checks it first.
	errmsg string
	err    error
}

func NewApp() *App {
	// Set up AWS
	//
	awsclients, err := aws.New()
	if err != nil {
		return &App{
			errmsg: "cannot set up AWS clients",
			err:    err,
		}
	}

	// Get postgres password
	//
	pgpwd := os.Getenv("PGPASSWORD")
	if pgpwd == "" {
		var err error
		pgpwd, err = awsclients.GetSecret(os.Getenv("FI_PG_SECRET_ID"))
		if err != nil {
			return &App{
				errmsg: "cannot get postgres password from secrets manager",
				err:    err,
			}
		}
	}

	// Open postgres connection
	//

	db, err := database.Open("pgx", database.GetDSN(pgpwd))
	if err != nil {
		return &App{
			errmsg: "cannot open connection to postgres",
			err:    err,
		}
	}

	// Grab some config from the environment
	//
	var maxmetrics int = 200000
	s := os.Getenv("MAX_METRICS")
	if s != "" {
		maxmetrics, err = strconv.Atoi(s)
		if err != nil {
			return &App{
				errmsg: "bad MAX_METRICS value",
				err:    err,
			}
		}
	}

	// Initialise our function's app
	//
	d, err := geodata.New(db, maxmetrics)
	if err != nil {
		return &App{
			errmsg: "cannot initialise geodata app",
			err:    err,
		}
	}

	return &App{
		d: d,
	}
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
	bbox := req.QueryStringParameters["bbox"]
	location := req.QueryStringParameters["location"]

	var radius int
	s := req.QueryStringParameters["radius"]
	if s != "" {
		var err error
		radius, err = strconv.Atoi(s)
		if err != nil {
			return errorResponse(http.StatusBadRequest, "malformed radius", err), nil
		}
	}
	geotypes := req.MultiValueQueryStringParameters["geotype"]
	// empty list means ALL
	rows := req.MultiValueQueryStringParameters["rows"]
	cols := req.MultiValueQueryStringParameters["cols"]

	body, err := app.d.Query(ctx, dataset, bbox, location, radius, geotypes, rows, cols)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "problem with query", err), nil
	}

	// headers
	headers := map[string]string{"Access-Control-Allow-Origin": "*"}

	response := &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       body,
		Headers:    headers,
	}

	return response, nil
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

func main() {
	app := NewApp()
	lambda.Start(app.Handler)
}
