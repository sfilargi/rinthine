package model

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

type User struct {
	UserId    []byte `json:"user_id" dynamodbav:"user_id_"`
	Handle    string `json:"handle" dynamodbav:"handle_"`
	Name      string `json:"name" dynamodbav:"name_"`
	Email     string `json:"email" dynamodbav:"email_"`
	Password  string `json:"password" dynamodbav:"password_"`
	AvatarUrl string `json:"avatar_url" dynamodbav:"avatar_url_"`
	Bio       string `json:"bio" dynamodbav:"bio_"`
	CreatedAt int64  `json:"created_at" dynamodbav:"created_at_"`
}

func UserPutTx(item *User) *types.TransactWriteItem {
	item.CreatedAt = time.Now().Unix()
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	return &types.TransactWriteItem{
		Put: &types.Put{
			Item:                i,
			TableName:           aws.String("core_users"),
			ConditionExpression: aws.String("attribute_not_exists(user_id_)"),
		},
	}
}

func UserPut(item *User) error {
	_, err := Db.TransactWriteItems(context.TODO(), &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			*UserIdxHandlePutTx(&UserIdxHandle{
				Handle: item.Handle,
				UserId: item.UserId,
			}),
			*UserPutTx(item),
		}})
	if err != nil {
		zap.S().Error(err.Error())
	}

	return err
}

func UserGet(pkey []byte) (*User, error) {
	result, err := Db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("core_users"),
		Key: map[string]types.AttributeValue{
			"user_id_": &types.AttributeValueMemberB{Value: pkey},
		},
	})
	if result.Item == nil || err != nil {
		return nil, err
	}

	var item User
	if err = attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		panic(err)
	}
	return &item, nil
}

func UserGetFromHandle(handle string) (*User, error) {
	fk, err := UserIdxHandleGet(handle)
	if fk == nil || err != nil {
		return nil, err
	}

	return UserGet(fk.UserId)
}
