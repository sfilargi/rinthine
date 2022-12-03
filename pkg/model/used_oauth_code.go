package model

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UsedOauthCode struct {
	Code string `json:"code" dynamodbav:"code_"`
}

func UsedOauthCodePut(item *UsedOauthCode) error {
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	_, err = Db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                i,
		TableName:           aws.String("core_used_oauth_codes"),
		ConditionExpression: aws.String("attribute_not_exists(code_)"),
	})

	return err
}
