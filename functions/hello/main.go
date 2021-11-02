package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	_ "github.com/lib/pq"
)

type App struct {
	aws *session.Session // AWS session
	sm  *secretsmanager.SecretsManager
	db  *sql.DB

	// It's hard to return usable errors from a lambda main function.
	// So any errors in main set err here and Handler checks it first.
	errmsg string
	err    error
}

func errorResponse(status int, msg string, err error) *events.APIGatewayProxyResponse {
	log.Printf("%s: %s", msg, err.Error())

	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}
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

func (app *App) Handler(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// return any init errors
	if app.err != nil {
		return errorResponse(http.StatusInternalServerError, app.errmsg, app.err), nil
	}

	// ping db before replying
	err := app.db.PingContext(ctx)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "cannot ping db", err), nil
	}

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
	return response, nil
}

func main() {
	app := NewApp()
	lambda.Start(app.Handler)
}
