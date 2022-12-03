package model

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	//"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AppIdxUserId struct {
	UserId []byte `json:"user_id" dynamodbav:"user_id_"`
	AppId  []byte `json:"app_id" dynamodbav:"app_id_"`
}

func AppIdxUserIdPut(item *AppIdxUserId) error {
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = Db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      i,
		TableName: aws.String("core_apps_idx_user_id"),
		ConditionExpression: aws.String(
			"attribute_not_exists(user_id_) AND attribute_not_exists(app_id_)"),
	})
	return err
}
