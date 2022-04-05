package domain

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type (
	Condition string

	WithCondition interface {
		SetName(name string) WithCondition
		Name() string
		Value() types.AttributeValue
		KeyCondition() string
	}

	WithSortKeyCondition interface {
		HasSortKey() bool
		WithCondition

		StarsWith(value interface{}) WithSortKeyCondition
		Equal(value interface{}) WithSortKeyCondition
		LessThan(value interface{}) WithSortKeyCondition
		LessThanOrEqual(value interface{}) WithSortKeyCondition
		GreaterThan(value interface{}) WithSortKeyCondition
		GreaterThanOrEqual(value interface{}) WithSortKeyCondition
		Between(start, end interface{}) WithSortKeyCondition

		SimpleCondition() bool
		StartValue() interface{}
		EndValue() interface{}
	}

	SqlExpression interface {
		SetIndex(indexName string) SqlExpression
		Where(condition WithCondition) SqlExpression
		AndWhere(keyCondition WithSortKeyCondition) SqlExpression
		ExpressionAttributeValues() map[string]types.AttributeValue
		IndexName() *string
		Update(keys ...WithCondition) SqlExpression
		UpdateExpression() *string
		AttributeNames() map[string]string

		Key() map[string]types.AttributeValue
		KeyCondition() *string

		SetItem(item interface{}) SqlExpression
		Names() map[string]types.AttributeValue
		Values() map[string]types.AttributeValue
	}
)
