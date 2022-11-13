package coreops

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Handle    string `json:"handle" dynamodbav:"handle"`
	Name      string `json:"name" dynamodbav:"name"`
	Email     string `json:"email" dynamodbav:"email"`
	Password  string `json:"password" dynamodbav:"password"`
	AvatarUrl string `json:"avatar_url" dynamodbav:"avatar_url"`
	CreatedAt int64  `json:"created_at" dynamodbav:"created_at"`
}

type UserToken struct {
	Token string `json:"token" dynamodbav:"token_"`
	User  string `json:"user" dynamodbav:"user"`
}

func VerifyPassword(user *User, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func RandomString(count int) (string, error) {
	buf := make([]byte, count)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(buf), nil
}

func CreateUserToken(handle string) (string, error) {

	token, err := RandomString(36)
	if err != nil {
		return "", err
	}
	userToken := UserToken{
		Token: token,
		User:  handle,
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())

	db := dynamodb.NewFromConfig(cfg)
	item, err := attributevalue.MarshalMap(userToken)
	if err != nil {
		return "", err
	}

	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String("core_user_tokens"),
		ConditionExpression: aws.String("attribute_not_exists(token_)"),
	})
	return token, err
}

func GetUser(handle string) (*User, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)

	result, err := db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"handle": &types.AttributeValueMemberS{Value: handle},
		},
		TableName: aws.String("core_users"),
	})
	if err != nil || result.Item == nil {
		return nil, err
	}

	var user User
	if err = attributevalue.UnmarshalMap(result.Item, &user); err != nil {
		return nil, fmt.Errorf("Internal Error")
	}
	return &user, nil
}

func CreateUser(user *User) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	db := dynamodb.NewFromConfig(cfg)
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("Failed to marshall request")
	}

	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String("core_users"),
		ConditionExpression: aws.String("attribute_not_exists(handle)"),
	})

	return err
}
