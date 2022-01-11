package dynamo_for_lambda

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

type (
	KeyCondition struct {
		name *string
		Val  interface{}
	}

	WithCondition interface {
		SetName(name string) WithCondition
		Name() string
		Value() types.AttributeValue
	}

	SortKeyCondition struct {
		exists bool
		name   *string
		Val    interface{}
	}

	WithSortKeyCondition interface {
		hasSortKey() bool
		WithCondition
	}

	Expression struct {
		indexName *string
		hashKey   *string
		rangeKey  *string

		expressions map[string]WithCondition
	}

	SqlExpression interface {
		SetIndex(indexName string) SqlExpression
		Where(condition WithCondition) SqlExpression
		AndWhere(keyCondition WithSortKeyCondition) SqlExpression
		Build() map[string]types.AttributeValue

		Key() map[string]types.AttributeValue
		KeyCondition() *string

		Names() map[string]types.AttributeValue
		Values() map[string]types.AttributeValue
	}
)

func NewSqlBuilder(config *Config) SqlExpression {
	return &Expression{
		hashKey:     aws.String(config.HashKeyName),
		rangeKey:    config.RangeKeyName,
		expressions: map[string]WithCondition{},
	}
}

/* KeyCondition */

func NewKeyCondition(name string, val interface{}) *KeyCondition {
	kc := &KeyCondition{Val: val}
	kc.SetName(name)

	return kc
}

func (k *KeyCondition) SetName(name string) WithCondition {
	k.name = aws.String(name)
	return k
}

func (k *KeyCondition) Name() string {
	return *k.name
}

func (k *KeyCondition) Value() types.AttributeValue {
	if k.Val == nil {
		log.Fatalf("%s has a nil value", k.Name())
	}

	return &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", k.Val)}
}

/* SortKeyCondition */

func NewSortKeyCondition(name string, val interface{}) *SortKeyCondition {
	ksc := &SortKeyCondition{Val: val}
	ksc.SetName(name)

	return ksc
}

func (k *SortKeyCondition) SetName(name string) WithCondition {
	k.name = aws.String(name)
	k.exists = true
	return k
}

func (k *SortKeyCondition) Name() string {
	return *k.name
}

func (k *SortKeyCondition) Value() types.AttributeValue {
	if k.Val == nil {
		log.Fatalf("%s has a nil value", k.Name())
	}

	return &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", k.Val)}
}

func (k *SortKeyCondition) hasSortKey() bool {
	return k.exists == true
}

/* Expression */

func (e *Expression) SetIndex(indexName string) SqlExpression {
	e.indexName = aws.String(indexName)
	return e
}

func (e *Expression) Where(condition WithCondition) SqlExpression {
	e.expressions["key"] = condition
	return e
}

func (e *Expression) AndWhere(keyCondition WithSortKeyCondition) SqlExpression {
	if keyCondition.hasSortKey() {
		e.expressions["sortKey"] = keyCondition
	}

	return e
}

func (e *Expression) Build() map[string]types.AttributeValue {
	built := map[string]types.AttributeValue{}
	return built
}

func (e *Expression) Key() map[string]types.AttributeValue {
	keys := map[string]types.AttributeValue{}

	for key, expr := range e.expressions {
		switch key {
		case "key", "sortKey":
			if expr.Name() != "" {
				keys[expr.Name()] = expr.Value()
			}
		}
	}

	return keys
}

func (e *Expression) KeyCondition() *string {
	return aws.String("")
}

func (e *Expression) Names() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{}
}

func (e *Expression) Values() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{}
}
