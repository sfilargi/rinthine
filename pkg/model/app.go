package model

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type App struct {
	AppId       []byte `json:"app_id" dynamodbav:"app_id_"`
	UserId      []byte `json:"user_id" dynamodbav:"user_id_"`
	Name        string `json:"name" dynamodbav:"name_"`
	Description string `json:"description" dynamodbav:"description_"`
	HomeUrl     string `json:"home_url" dynamodbav:"home_url_"`
	RedirectUrl string `json:"redirect_url" dynamodbav:"redirect_url_"`
	Password    string `json:"-" dynamodbav:"password_"`
	CreatedAt   int64  `json:"created_at" dynamodbav:"created_at_"`
}

func AppPutTx(item *App) *types.TransactWriteItem {
	item.CreatedAt = time.Now().Unix()
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	return &types.TransactWriteItem{
		Put: &types.Put{
			Item:                i,
			TableName:           aws.String("core_apps"),
			ConditionExpression: aws.String("attribute_not_exists(app_id_)"),
		},
	}
}

func AppPut(item *App) error {
	// We don't need transaction for App indexes
	err := AppIdxUserIdPut(&AppIdxUserId{
		UserId: item.UserId,
		AppId:  item.AppId,
	})
	if err != nil {
		return err
	}

	item.CreatedAt = time.Now().Unix()
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	_, err = Db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                i,
		TableName:           aws.String("core_apps"),
		ConditionExpression: aws.String("attribute_not_exists(app_id_)"),
	})

	return err
}

func AppGet(app_id []byte) (*App, error) {
	result, err := Db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("core_apps"),
		Key: map[string]types.AttributeValue{
			"app_id_": &types.AttributeValueMemberB{Value: app_id},
		},
	})
	if result.Item == nil || err != nil {
		return nil, err
	}

	var item App
	if err = attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		panic(err)
	}
	return &item, nil
}
