package model

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	//"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

type OauthCodeIdxAppId struct {
	AppId []byte `json:"app_id" dynamodbav:"app_id_"`
	Code  string `json:"code" dynamodbav:"code_"`
}

func OauthCodeIdxAppIdPut(item *OauthCodeIdxAppId) error {
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	_, err = Db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      i,
		TableName: aws.String("core_oauth_codes_idx_app_id"),
		ConditionExpression: aws.String(
			"attribute_not_exists(app_id_) AND attribute_not_exists(code_)"),
	})
	if err != nil {
		zap.S().Error(err.Error())
	}
	return err
}
