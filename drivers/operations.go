package drivers

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
)

func (d *DynamoClient) Get(expression domain.SqlExpression, target interface{}) error {
	output, err := d.Client.GetItem(d.Ctx, &dynamodb.GetItemInput{
		TableName: d.TableName,
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

func (d *DynamoClient) Query(expression domain.SqlExpression, target interface{}) error {
	output, err := d.Client.Query(d.Ctx, &dynamodb.QueryInput{
		TableName:                 d.TableName,
		KeyConditionExpression:    expression.KeyCondition(),
		ExpressionAttributeValues: expression.ExpressionAttributeValues(),
		IndexName:                 expression.IndexName(),
	})

	if err != nil {
		return fmt.Errorf("query: %v", err)
	}

	err = attributevalue.UnmarshalListOfMaps(output.Items, target)
	if err != nil {
		return fmt.Errorf("UnmarshalMap: %v", err)
	}

	return nil
}

func (d *DynamoClient) Put(item domain.SqlExpression, result interface{}) error {
	_, err := d.Client.PutItem(d.Ctx, &dynamodb.PutItemInput{
		Item:      item.Values(),
		TableName: d.TableName,
	})
	if err != nil {
		return fmt.Errorf("put item: %v", err)
	}

	err = attributevalue.UnmarshalMap(item.Values(), result)
	if err != nil {
		return fmt.Errorf("UnmarshalMap: %v", err)
	}

	return nil
}

func (d *DynamoClient) Update(expression domain.SqlExpression, result interface{}) error {
	out, err := d.Client.UpdateItem(d.Ctx, &dynamodb.UpdateItemInput{
		TableName:                 d.TableName,
		Key:                       expression.Key(),
		UpdateExpression:          expression.UpdateExpression(),
		ExpressionAttributeValues: expression.ExpressionAttributeValues(),
		ExpressionAttributeNames:  expression.AttributeNames(),
	})

	if err != nil {
		return fmt.Errorf("update item: %v", err)
	}

	err = attributevalue.UnmarshalMap(out.Attributes, result)
	if err != nil {
		return fmt.Errorf("UnmarshalMap: %v", err)
	}

	return nil
}

func (d *DynamoClient) Delete(expression domain.SqlExpression) error {
	out, err := d.Client.DeleteItem(d.Ctx, &dynamodb.DeleteItemInput{
		TableName: d.TableName,
		Key:       expression.Key(),
	})

	if err != nil {
		return fmt.Errorf("delete item: %v", err)
	}

	err = attributevalue.UnmarshalMap(out.Attributes, nil)
	if err != nil {
		return fmt.Errorf("UnmarshalMap: %v", err)
	}

	return nil
}
