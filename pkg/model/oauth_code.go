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

type OauthCode struct {
	Code      string `json:"code" dynamodbav:"code_"`
	AppId     []byte `json:"app_id" dynamodbav:"app_id_"`
	UserId    []byte `json:"user_id" dynamodbav:"user_id_"`
	CreatedAt int64  `json:"created_at" dynamodbav:"created_at_"`
}

func OauthCodePutTx(item *OauthCode) *types.TransactWriteItem {
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	return &types.TransactWriteItem{
		Put: &types.Put{
			Item:                i,
			TableName:           aws.String("core_oath_codes"),
			ConditionExpression: aws.String("attribute_not_exists(code_)"),
		},
	}
}

func OauthCodePut(item *OauthCode) error {
	// We don't need transaction for OauthCode indexes
	err := OauthCodeIdxUserIdPut(&OauthCodeIdxUserId{
		UserId: item.UserId,
		Code:   item.Code,
	})
	if err != nil {
		return err
	}

	err = OauthCodeIdxAppIdPut(&OauthCodeIdxAppId{
		AppId: item.AppId,
		Code:  item.Code,
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
		TableName:           aws.String("core_oauth_codes"),
		ConditionExpression: aws.String("attribute_not_exists(code_)"),
	})
	if err != nil {
		zap.S().Error(err.Error())
	}

	return err
}

func OauthCodeGet(code string) (*OauthCode, error) {
	result, err := Db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("core_oauth_codes"),
		Key: map[string]types.AttributeValue{
			"code_": &types.AttributeValueMemberS{Value: code},
		},
	})
	if result.Item == nil || err != nil {
		if err != nil {
			zap.S().Error(err.Error())
		}
		return nil, err
	}

	var item OauthCode
	if err = attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		panic(err)
	}
	return &item, nil
}
