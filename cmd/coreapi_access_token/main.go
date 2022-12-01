package main

import (
	"context"

	"github.com/rinthine/pkg/coreops"
	"github.com/rinthine/pkg/misc"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	form, err := misc.ParseForm(e)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	// grant_type
	code := form["code"][0] // let it panic
	// redirect_uri

	authstring, _ := e.Headers["authorization"]
	appname, password, err := coreops.BasicCredentials(authstring)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	app, err := coreops.GetApp(appname)
	if app == nil || err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	if app.Password != password {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid password",
		}, nil
	}

	token, err := coreops.CreateAppToken(app, code)
	if err != nil {
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
