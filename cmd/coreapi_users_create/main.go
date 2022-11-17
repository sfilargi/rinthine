package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rinthine/pkg/coreops"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type UserCreateRequest struct {
	Handle   string `json:"handle"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Handle must be 6+ characters
// Password must not be empty
// Email and Password "validation" is left to front end. We don't care
// We don't sanitize!
func verifyInput(req *UserCreateRequest) error {
	if len(req.Handle) < 6 {
		return fmt.Errorf("handle must be 6+ characters long")
	}
	if len(req.Password) == 0 {
		return fmt.Errorf("password cannot be empty")
	}
	if len(req.Name) == 0 {
		return fmt.Errorf("name canot be empty")
	}
	return nil
}

func HandleRequest(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var req UserCreateRequest
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

	existing_user, err := coreops.GetUser(req.Handle)
	if existing_user != nil || err != nil {
		body := "handle taken"
		if err != nil {
			body = err.Error()
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       body,
		}, nil
	}

	user := coreops.User{
		UserId:    coreops.GenId(),
		Handle:    req.Handle,
		Name:      req.Name,
		Email:     req.Email,
		Password:  coreops.HashPassword(req.Password),
		CreatedAt: time.Now().Unix(),
	}
	if err := coreops.CreateUser(&user); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	token, err := coreops.CreateUserToken(user.Handle)
	if err != nil {
		// hmm, do we return success or failure here?
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       token,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

