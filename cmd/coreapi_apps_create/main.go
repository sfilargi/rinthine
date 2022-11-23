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
	Name     string `json:"name"`
	Description string `json:"description"`
	HomeUrl string `json:"home_url"`
	RedirectUrl string `json:"redirect_url"`
}

// Handle must be 6+ characters
// Password must not be empty
// Email and Password "validation" is left to front end. We don't care
// We don't sanitize!
func verifyInput(req *AppCreateRequest) error {
	if len(req.Name) < 6 {
		return fmt.Errorf("name must be 6+ characters long")
	}
	if len(req.RedirectUrl) == 0 {
		return fmt.Errorf("redirect_url canot be empty")
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
	if err := verifyInput(&req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
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

	data, err := json.Marshal(user)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(data),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}