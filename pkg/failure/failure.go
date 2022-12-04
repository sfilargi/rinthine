package failure

import (
	"github.com/aws/aws-lambda-go/events"
)

type Code int64

const (
	InternalError Code = iota
	InvalidName
	NoRedirectUrl
	HandleTaken
	InvalidHandle
	InvalidPassword
)

type Error interface {
	ToL() (events.APIGatewayProxyResponse, error)
}

func (f Code) ToL() (events.APIGatewayProxyResponse, error) {
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
	case HandleTaken:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "{\"message\": \"Handle taken\"}",
		}, nil
	case InvalidHandle:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "{\"message\": \"Invalid Handle\"}",
		}, nil
	case InvalidPassword:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "{\"message\": \"Invalid Password\"}",
		}, nil
	}

	panic("oops")
}
