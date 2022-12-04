package model

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	//"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

type OauthCodeIdxUserId struct {
	UserId []byte `json:"user_id" dynamodbav:"user_id_"`
	Code   string `json:"code" dynamodbav:"code_"`
}

func OauthCodeIdxUserIdPut(item *OauthCodeIdxUserId) error {
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = Db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      i,
		TableName: aws.String("core_oauth_codes_idx_user_id"),
		ConditionExpression: aws.String(
			"attribute_not_exists(user_id_) AND attribute_not_exists(code_)"),
	})
	if err != nil {
		zap.S().Error(err.Error())
	}
	return err
}
