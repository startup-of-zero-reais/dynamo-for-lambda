package domain

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type (
	WithCondition interface {
		SetName(name string) WithCondition
		Name() string
		Value() types.AttributeValue
	}

	WithSortKeyCondition interface {
		HasSortKey() bool
		WithCondition
	}

	SqlExpression interface {
		SetIndex(indexName string) SqlExpression
		Where(condition WithCondition) SqlExpression
		AndWhere(keyCondition WithSortKeyCondition) SqlExpression
		Build() map[string]types.AttributeValue

		Key() map[string]types.AttributeValue
		KeyCondition() *string

		SetItem(item interface{}) SqlExpression
		Names() map[string]types.AttributeValue
		Values() map[string]types.AttributeValue
	}
)
