package coreops

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Handle    string `json:"handle" dynamodbav:"handle_"`
	Name      string `json:"name" dynamodbav:"name_"`
	Email     string `json:"email" dynamodbav:"email_"`
	Password  string `json:"-" dynamodbav:"password_"`
	AvatarUrl string `json:"avatar_url" dynamodbav:"avatar_url_"`
	CreatedAt int64  `json:"created_at" dynamodbav:"created_at_"`
}

type UserToken struct {
	Token string `dynamodbav:"token_"`
	User  string `dynamodbav:"user_"`
	App   string `dynamodbav:"app_"`
}

type App struct {
	Name        string `json:"name" dynamodbav:"name_"`
	User        string `json:"user" dynamodbav:"user_"`
	Description string `json:"description" dynamodbav:"description_"`
	HomeUrl     string `json:"home_url" dynamodbav:"home_url_"`
	RedirectUrl string `json:"redirect_url" dynamodbav:"redirect_url_"`
	Password    string `json:"-" dynamodbav:"password_"`
}

type UserApp struct {
	User string `json:"user" dynamodbav:"user_"`
	App  string `json:"app" dynamodbav:"app_"`
}

func SecureRandUint32() uint32 {
	buffer := make([]byte, 4)
	n, err := rand.Read(buffer)
	if n != 4 || err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint32(buffer)
}

func GenId() []byte {
	var start uint64 = 1641024000 // 1/1/2022 00:00:00
	v := ((uint64(time.Now().Unix()) - start) << 32) + uint64(SecureRandUint32())
	id := make([]byte, 8)
	binary.LittleEndian.PutUint64(id, v)
	return id
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

func RandomString(count int) string {
	buf := make([]byte, count)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}

func CreateUserToken(user string) (string, error) {

	userToken := UserToken{
		Token: RandomString(36),
		User:  user,
	}

	cfg, _ := config.LoadDefaultConfig(context.TODO())

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
	return userToken.Token, err
}

func GetUser(handle string) (*User, error) {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)

	result, err := db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("core_users"),
		Key: map[string]types.AttributeValue{
			"handle_": &types.AttributeValueMemberS{Value: handle},
		},
	})
	if result.Item == nil || err != nil {
		return nil, err
	}

	var user User
	if err = attributevalue.UnmarshalMap(result.Item, &user); err != nil {
		return nil, fmt.Errorf("Internal Error")
	}
	return &user, nil
}

func CreateUser(user *User) error {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("Failed to marshall request")
	}

	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String("core_users"),
		ConditionExpression: aws.String("attribute_not_exists(handle_)"),
	})

	return err
}

func BearerToken(authstring string) (string, error) {
	ss := strings.SplitN(authstring, " ", 2)
	if len(ss) != 2 {
		return "", fmt.Errorf("Couldn't parse Authorization header %s, %v, %d",
			authstring, ss, len(ss))
	}
	if !strings.EqualFold(ss[0], "Bearer") {
		return "", fmt.Errorf("Not a Bearer Authorization")
	}
	return ss[1], nil
}

func GetUserToken(token string) (*UserToken, error) {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)

	result, err := db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("core_user_tokens"),
		Key: map[string]types.AttributeValue{
			"token_": &types.AttributeValueMemberS{Value: token},
		},
	})
	if err != nil || result.Item == nil {
		return nil, err
	}

	var userToken UserToken
	if err = attributevalue.UnmarshalMap(result.Item, &userToken); err != nil {
		return nil, fmt.Errorf("Internal Error")
	}
	return &userToken, nil
}

func BearerAuthenticate(token string) (*User, error) {
	userToken, err := GetUserToken(token)
	if userToken == nil || err != nil {
		return nil, err
	}

	user, err := GetUser(userToken.User)
	return user, nil
}

func GetApp(name string) (*App, error) {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)

	result, err := db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("core_apps"),
		Key: map[string]types.AttributeValue{
			"name_": &types.AttributeValueMemberS{Value: name},
		},
	})
	if err != nil || result.Item == nil {
		return nil, err
	}

	var app App
	if err = attributevalue.UnmarshalMap(result.Item, &app); err != nil {
		return nil, fmt.Errorf("Internal Error")
	}
	return &app, nil
}

func GetUserApp(user string) (*App, error) {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)

	result, err := db.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("core_user_apps"),
		KeyConditionExpression: aws.String("user_ = :user"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user": &types.AttributeValueMemberS{Value: user},
		},
		Limit: aws.Int32(1),
	})
	if err != nil || len(result.Items) == 0 {
		return nil, err
	}

	var userApp UserApp
	if err = attributevalue.UnmarshalMap(result.Items[0], &userApp); err != nil {
		return nil, fmt.Errorf("Internal Error")
	}

	return GetApp(userApp.App)
}

func CreateUserApp(user, app string) error {

	old, err := GetUserApp(user)
	if old != nil {
		return fmt.Errorf("Only one app allowed")
	}
	if err != nil {
		return err
	}

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)
	item, err := attributevalue.MarshalMap(UserApp{
		User: user,
		App:  app,
	})
	if err != nil {
		return fmt.Errorf("Failed to marshall request")
	}

	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String("core_user_apps"),
		ConditionExpression: aws.String("attribute_not_exists(user_)"),
	})

	return err
}

func CreateApp(app *App) error {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)
	item, err := attributevalue.MarshalMap(app)
	if err != nil {
		return fmt.Errorf("Failed to marshall request")
	}

	err = CreateUserApp(app.User, app.Name)
	if err != nil {
		return err
	}

	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String("core_apps"),
		ConditionExpression: aws.String("attribute_not_exists(name_)"),
	})

	return err
}
