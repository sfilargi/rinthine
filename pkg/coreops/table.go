package coreops

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Table struct {
	name          string
	partition_key string
	sort_key      string
}

func CreateTable(name, partition_key, sort_key string) *Table {
	return &Table{
		name:          name,
		partition_key: partition_key,
		sort_key:      sort_key,
	}
}

func (t *Table) CreateItem(item interface{}) error {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)

	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      i,
		TableName: aws.String(t.name),
		ConditionExpression: aws.String(
			fmt.Sprintf("attribute_not_exists(%s)", t.partition_key)),
	})

	return err
}

func (t *Table) AddItem(item interface{}) error {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)

	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      i,
		TableName: aws.String(t.name),
	})

	return err
}

func (t *Table) GetItemS(key string, item interface{}) error {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	db := dynamodb.NewFromConfig(cfg)

	result, err := db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(t.name),
		Key: map[string]types.AttributeValue{
			t.partition_key: &types.AttributeValueMemberS{Value: key},
		},
	})
	if result.Item == nil || err != nil {
		return err
	}

	if err = attributevalue.UnmarshalMap(result.Item, item); err != nil {
		return fmt.Errorf("Internal Error")
	}
	return nil
}
