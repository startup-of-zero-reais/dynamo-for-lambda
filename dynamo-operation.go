package dynamo_for_lambda

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type (
	DynamoSQL interface {
		Get(expression SqlExpression, target interface{}) error
		Put(item interface{}, result interface{}) error
		Update(expression Expression, item interface{}, result interface{}) error
		Delete(expression SqlExpression) error
	}
)

func (d *DynamoClient) Get(expression SqlExpression, target interface{}) error {
	output, err := d.Client.GetItem(d.Ctx, &dynamodb.GetItemInput{
		TableName: d.table,
		Key:       expression.Key(),
	})
	if err != nil {
		return fmt.Errorf("get item: %v", err)
	}

	err = attributevalue.UnmarshalMap(output.Item, target)
	if err != nil {
		return fmt.Errorf("UnmarshalMap: %v", err)
	}

	return nil
}

func (d *DynamoClient) Put(item interface{}, result interface{}) error {
	return nil
}

func (d *DynamoClient) Update(expression Expression, item interface{}, result interface{}) error {
	return nil
}

func (d *DynamoClient) Delete(expression SqlExpression) error {
	return nil
}
