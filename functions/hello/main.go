package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
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
	lambda.Start(handler)
}
