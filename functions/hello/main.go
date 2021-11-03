package main

import (
	"context"
	"database/sql"
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

func errorResponse(status int, logmsg string, err error, usermsg string) *events.APIGatewayProxyResponse {
	if err == nil {
		log.Println(logmsg)
	} else {
		log.Printf("%s: %s", logmsg, err.Error())
	}
	if usermsg == "" {
		usermsg = logmsg
	}

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusText(status),
		Message: usermsg,
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
	return errorResponse(http.StatusBadRequest, msg, nil, "")
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

func (app *App) Handler(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// return any init errors
	if app.err != nil {
		return errorResponse(http.StatusInternalServerError, app.errmsg, app.err, "problem setting up database connection"), nil
	}

	// grab parameters from path and query string
	//
	dataset := req.PathParameters["dataset"]
	if dataset == "" {
		return clientResponse("missing dataset path parameter"), nil
	}
	rows := gatherTokens(req.MultiValueQueryStringParameters["rows"])
	if len(rows) == 0 {
		return clientResponse("missing rows query parameter"), nil
	}
	cols := gatherTokens(req.MultiValueQueryStringParameters["cols"])
	if len(cols) == 0 {
		return clientResponse("missing cols query parameter"), nil
	}

	// ping db before replying
	err := app.db.PingContext(ctx)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "cannot ping db", err, "problem pinging database"), nil
	}

	response := &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("dataset=%s\nrows=%v\ncols=%v\n", dataset, rows, cols),
	}
	/*
		buf, err := json.MarshalIndent(*req, "", "    ")
		if err != nil {
			response := &events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       http.StatusText(http.StatusInternalServerError),
			}
			return response, nil
		}
		response := &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(buf),
		}
	*/
	return response, nil
}

func main() {
	app := NewApp()
	lambda.Start(app.Handler)
}
