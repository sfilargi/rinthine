package model

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Token struct {
	Token     string `json:"token" dynamodbav:"token_"`
	UserId    []byte `json:"user_id" dynamodbav:"user_id_"`
	AppId     []byte `json:"app_id" dynamodbav:"app_id_"`
	CreatedAt int64  `json:"created_at" dynamodbav:"created_at_"`
}

func TokenPutTx(item *Token) *types.TransactWriteItem {
	item.CreatedAt = time.Now().Unix()
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	return &types.TransactWriteItem{
		Put: &types.Put{
			Item:                i,
			TableName:           aws.String("core_tokens"),
			ConditionExpression: aws.String("attribute_not_exists(token_)"),
		},
	}
}

func TokenPut(item *Token) error {
	// We don't need transaction for Token indexes
	err := TokenIdxUserIdPut(&TokenIdxUserId{
		UserId: item.UserId,
		Token:  item.Token,
	})
	if err != nil {
		return err
	}

	if item.AppId != nil {
		err := TokenIdxAppIdPut(&TokenIdxAppId{
			AppId: item.AppId,
			Token: item.Token,
		})
		if err != nil {
			return err
		}
	}

	item.CreatedAt = time.Now().Unix()
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	_, err = Db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:                i,
		TableName:           aws.String("core_tokens"),
		ConditionExpression: aws.String("attribute_not_exists(token_)"),
	})

	return err
}

func TokenGet(token string) (*Token, error) {
	result, err := Db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("core_tokens"),
		Key: map[string]types.AttributeValue{
			"token_": &types.AttributeValueMemberS{Value: token},
		},
	})
	if result.Item == nil || err != nil {
		return nil, err
	}

	var item Token
	if err = attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		panic(err)
	}
	return &item, nil
}
