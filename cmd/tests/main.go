package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rinthine/pkg/coreops"

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

type Test struct {
	A     string `dynamodbav:"a"`
	Id    []byte `dynamodbav:"id_"`
	Order int
}

func TestPut(item *Test) error {
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	_, err = Db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      i,
		TableName: aws.String("genid_test"),
	})
	if err != nil {
		panic(err)
	}

	return err
}

func TestGetAll() {

	r, err := Db.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("genid_test"),
		KeyConditionExpression: aws.String("a = :a"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":a": &types.AttributeValueMemberS{Value: "a"},
		},
	})
	if err != nil {
		panic(err)
	}
	for _, item := range r.Items {
		var i Test
		if err = attributevalue.UnmarshalMap(item, &i); err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", i)
	}

}

func printBin(bin []byte) {
	for _, b := range bin {
		fmt.Printf("%03d ", b)
	}
	fmt.Println()
}

func main() {
	for i := 0; i < 10; i++ {
		TestPut(&Test{
			A:     "a",
			Id:    coreops.GenId(),
			Order: i,
		})
		time.Sleep(1 * time.Second)
	}
	TestGetAll()
}
