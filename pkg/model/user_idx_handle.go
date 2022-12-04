package model

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

type UserIdxHandle struct {
	Handle string `json:"handle" dynamodbav:"handle_"`
	UserId []byte `json:"user_id" dynamodbav:"user_id_"`
}

func UserIdxHandlePutTx(item *UserIdxHandle) *types.TransactWriteItem {
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	return &types.TransactWriteItem{
		Put: &types.Put{
			Item:      i,
			TableName: aws.String("core_users_idx_handle"),
			ConditionExpression: aws.String(
				fmt.Sprintf(
					"attribute_not_exists(%s)", "handle_")),
		},
	}
}

func UserIdxHandlePut(item *UserIdxHandle) error {
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	_, err = Db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                i,
		TableName:           aws.String("core_users_idx_handle"),
		ConditionExpression: aws.String("attribute_not_exists(handle_)"),
	})
	if err != nil {
		zap.S().Error(err.Error())
	}
	return err
}

func UserIdxHandleGet(pkey string) (*UserIdxHandle, error) {
	result, err := Db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("core_users_idx_handle"),
		Key: map[string]types.AttributeValue{
			"handle_": &types.AttributeValueMemberS{Value: pkey},
		},
	})
	if result.Item == nil || err != nil {
		if err != nil {
			zap.S().Error(err.Error())
		}
		return nil, err
	}

	var item UserIdxHandle
	if err = attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		panic(err)
	}
	return &item, nil
}
