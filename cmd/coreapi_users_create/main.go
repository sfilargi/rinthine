package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Handle    string  `json:"handle" dynamodbav:"handle"`
	Name      *string `json:"name,omitempty" dynamodbav:"name"`
	Email     *string `json:"email,omitempty" dynamodbav:"email"`
	Password  *string `json:"password,omitempty" dynamodbav:"password"`
	AvatarUrl *string `json:"avatar_url,omitempty" dynamodbav:"avatar_url"`
	CreatedAt int64   `json:"created_at" dynamodbav:"created_at"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Handle must be 6+ characters
// Password must not be empty
// Email and Password "validation" is left to front end. We don't care
// We don't sanitize!
func verifyInput(user *User) error {
	if len(user.Handle) < 6 {
		return fmt.Errorf("handle must be 6+ characters long")
	}
	if user.Password == nil || len(*user.Password) == 0 {
		return fmt.Errorf("password cannot be empty")
	}
	if user.Name == nil || len(*user.Name) == 0 {
		return fmt.Errorf("name canot be empty")
	}
	if user.AvatarUrl != nil {
		return fmt.Errorf("avatar_url cannot be edited manually")
	}

	// hash the password
	var err error
	*user.Password, err = hashPassword(*user.Password)
	if err != nil {
		return fmt.Errorf("internal error trying to hash the password")
	}

	user.CreatedAt = time.Now().Unix()
	return nil
}

func createUser(user *User) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	db := dynamodb.NewFromConfig(cfg)
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("Failed to marshall request")
	}

	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String("users"),
		ConditionExpression: aws.String("attribute_not_exists(handle)"),
	})
	return err
}

func HandleRequest(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var user User
	if err := json.Unmarshal([]byte(e.Body), &user); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	if err := verifyInput(&user); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	if err := createUser(&user); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "done",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
