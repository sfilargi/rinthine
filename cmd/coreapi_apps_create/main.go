package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rinthine/pkg/coreops"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type AppCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HomeUrl     string `json:"home_url"`
	RedirectUrl string `json:"redirect_url"`
}

// Handle must be 6+ characters
// Password must not be empty
// Email and Password "validation" is left to front end. We don't care
// We don't sanitize!
func verifyInput(req *AppCreateRequest) Failure {
	if len(req.Name) < 6 {
		return InvalidName
	}
	if len(req.RedirectUrl) == 0 {
		return NoRedirectUrl
	}
	return nil
}

func HandleRequest(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var req AppCreateRequest
	if err := json.Unmarshal([]byte(e.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}
	if f := verifyInput(&req); f != nil {
		return f
	}

	authstring, _ := e.Headers["authorization"]
	token, err := coreops.BearerToken(authstring)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	user, err := coreops.BearerAuthenticate(token)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	app := coreops.App{
		Name:        req.Name,
		User:        user.Handle,
		Description: req.Description,
		HomeUrl:     req.HomeUrl,
		RedirectUrl: req.RedirectUrl,
		Password:    coreops.RandomString(36),
	}

	err = coreops.CreateApp(&app)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       app.Password,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
