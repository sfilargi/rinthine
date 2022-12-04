package model

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	//"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

type TokenIdxUserId struct {
	UserId []byte `json:"user_id" dynamodbav:"user_id_"`
	Token  string `json:"token" dynamodbav:"token_"`
}

func TokenIdxUserIdPut(item *TokenIdxUserId) error {
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	_, err = Db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      i,
		TableName: aws.String("core_tokens_idx_user_id"),
		ConditionExpression: aws.String(
			"attribute_not_exists(user_id_) AND attribute_not_exists(token_)"),
	})
	if err != nil {
		zap.S().Error(err.Error())
	}
	return err
}
