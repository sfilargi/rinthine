package coreops

import (
	"github.com/aws/aws-lambda-go/events"
)

type FailCode int64

const (
	InternalError FailCode = iota
	InvalidName
	NoRedirectUrl
)

type Failure interface {
	ToL() (events.APIGatewayProxyResponse, error)
}

func (f FailCode) ToL() (events.APIGatewayProxyResponse, error) {
	switch f {
	case InternalError:
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "{\"message\": \"Internal Error\"}",
		}, nil
	case InvalidName:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "{\"message\": \"Invalid Name\"}",
		}, nil
	case NoRedirectUrl:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "{\"message\": \"RedirectUrl can't be empty\"}",
		}, nil
	}
}
