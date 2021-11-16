package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/demo"
)

type App struct {
	d *demo.Demo

	// It's hard to return usable errors from a lambda main function.
	// So any errors in main set err here and Handler checks it first.
	errmsg string
	err    error
}

func NewApp() *App {
	d, err := demo.New(os.Getenv("PGPASSWORD"))
	app := &App{
		d: d,
	}
	if err != nil {
		app.err = err
		app.errmsg = err.Error()
	}
	return app
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
	rows := req.MultiValueQueryStringParameters["rows"]
	cols := req.MultiValueQueryStringParameters["cols"]

	body, err := app.d.Query(ctx, dataset, rows, cols)
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
