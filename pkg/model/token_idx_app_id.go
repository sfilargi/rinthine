package model

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	//"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TokenIdxAppId struct {
	AppId []byte `json:"app_id" dynamodbav:"app_id_"`
	Token string `json:"token" dynamodbav:"token_"`
}

func TokenIdxAppIdPut(item *TokenIdxAppId) error {
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = Db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      i,
		TableName: aws.String("core_tokens_idx_app_id"),
		ConditionExpression: aws.String(
			"attribute_not_exists(app_id_) AND attribute_not_exists(token_)"),
	})
	return err
}
