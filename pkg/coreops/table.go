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

func configMust() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	return cfg
}

var Db = dynamodb.NewFromConfig(configMust())

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
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = Db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      i,
		TableName: aws.String(t.name),
		ConditionExpression: aws.String(
			fmt.Sprintf("attribute_not_exists(%s)", t.partition_key)),
	})

	return err
}

func (t *Table) AddItem(item interface{}) error {
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = Db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      i,
		TableName: aws.String(t.name),
	})

	return err
}

func (t *Table) GetItemS(key string, item interface{}) error {
	result, err := Db.GetItem(context.TODO(), &dynamodb.GetItemInput{
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
