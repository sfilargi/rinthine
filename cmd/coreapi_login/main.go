package main

import (
	"context"
	"encoding/json"

	"github.com/rinthine/pkg/coreops"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LoginReq struct {
	Handle   string `json:"handle"`
	Password string `json:"password"`
}

func HandleRequest(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var req LoginReq
	if err := json.Unmarshal([]byte(e.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	user, err := coreops.GetUserByHandle(req.Handle)
	if user == nil || err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	if !coreops.VerifyPassword(user, req.Password) {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid password",
		}, nil
	}

	token, err := coreops.CreateUserToken(user.UserId)
	if err != nil {
		// hmm, do we return success or failure here?
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
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
