package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"

	"github.com/rinthine/pkg/coreops"
	"github.com/rinthine/pkg/misc"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//go:embed index.html
var index_data []byte
var index = template.Must(template.New("index").Parse(string(index_data)))

func HandlePost(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//data, _ := json.Marshal(&e)

	//client_id, _ := e.QueryStringParameters["client_id"]

	form, err := misc.ParseForm(e)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	username := form["username"][0] // Let it panic
	password := form["password"][0] // Let it panic
	clientId := e.QueryStringParameters["client_id"]

	app, code, err := coreops.Authorize(clientId, username, password)
	if app == nil || code == "" || err != nil {
		var body bytes.Buffer
		index.Execute(&body, nil)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "text/html",
			},
			Body: body.String(),
		}, nil
	}

	location := fmt.Sprintf("%s?code=%s", app.RedirectUrl, code)
	return events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"Content-Type": "text/html",
			"Location":     location,
		},
	}, nil
}

func HandleGet(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body bytes.Buffer
	index.Execute(&body, nil)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: body.String(),
	}, nil
}

func HandleRequest(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if e.HTTPMethod == "GET" {
		return HandleGet(ctx, e)
	} else if e.HTTPMethod == "POST" {
		return HandlePost(ctx, e)
	} else {
		panic(fmt.Sprintf("Unexpected HttpMethod: %s", e.HTTPMethod))
	}

	panic("oops")
}

func main() {
	misc.SetupZap()
	lambda.Start(HandleRequest)
}
