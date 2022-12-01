package misc

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

func ParseURLEncoded(s string) (map[string][]string, error) {
	return url.ParseQuery(s)
}

func ParseForm(e events.APIGatewayProxyRequest) (map[string][]string, error) {
	body := e.Body
	if e.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			panic(err)
		}
		body = string(decoded)
	}

	// This shit is case sensitive...
	if e.Headers["Content-Type"] == "application/x-www-form-urlencoded" {
		return ParseURLEncoded(body)
	} else {
		panic(fmt.Sprintf("Can't handle Content-Type: %+v", e.Headers))
	}
}
