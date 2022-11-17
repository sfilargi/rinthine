package coreops

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/crypto/bcrypt"
)

type UserHandle struct {
	Handle string `dynamodbav:"handle"`
	UserId []byte `dynamodbav:"userid"`
}

type UserId struct {
	UserId []byte `dynamodbav:"userid"`
	Handle string `dynamodbav:"handle"`
}

type User struct {
	UserId    []byte `dynamodbav:"userid"`
	Handle    string `dynamodbav:"handle"`
	Name      string `dynamodbav:"name"`
	Email     string `dynamodbav:"email"`
	Password  string `dynamodbav:"password"`
	AvatarUrl string `dynamodbav:"avatar_url"`
	CreatedAt int64  `dynamodbav:"created_at"`
}

type UserToken struct {
	Token string `dynamodbav:"token_"`
	User  string `dynamodbav:"user"`
	App   string `dynamodbav:"app"`
}

func SecureRandUint32() uint32 {
	buffer := make([]byte, 4);
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
	return token, err
}

func GetUserId(handle string) ([]byte, error) {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)
	result, err := db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"handle": &types.AttributeValueMemberS{Value: handle},
		},
		TableName: aws.String("core_userhandles"),
	})
	if err != nil || result.Item == nil {
		return nil, err
	}

	var userid UserId
	if err = attributevalue.UnmarshalMap(result.Item, &userid); err != nil {
		return nil, fmt.Errorf("Internal Error")
	}
	return userid.UserId, nil
	
}

func GetUser(handle string) (*User, error) {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)

	userid, err := GetUserId(handle)
	if userid == nil || err != nil {
		return nil, err
	}

	result, err := db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"userid": &types.AttributeValueMemberB{Value: userid},
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


// Some race conditions here, but we don't care much
func ReserveHandle(user *User) error {
	handle := UserHandle{
		Handle: user.Handle,
		UserId: user.UserId,
	}

	userid := UserId{
		UserId: user.UserId,
		Handle: user.Handle,
	}

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)

	handleitem, err := attributevalue.MarshalMap(handle)
	if err != nil {
		return fmt.Errorf("Failed to marshall request")
	}
	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                handleitem,
		TableName:           aws.String("core_userhandles"),
		ConditionExpression: aws.String("attribute_not_exists(handle)"),
	})

	useriditem, err := attributevalue.MarshalMap(userid)
	if err != nil {
		return fmt.Errorf("Failed to marshall request")
	}
	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                useriditem,
		TableName:           aws.String("core_userids"),
		ConditionExpression: aws.String("attribute_not_exists(userid)"),
	})

	return err
}

func CreateUser(user *User) error {

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("Failed to marshall request")
	}

	err = ReserveHandle(user)
	if err != nil {
		return err
	}

	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String("core_users"),
	})

	return err
}
