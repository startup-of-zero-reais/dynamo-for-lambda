package expressions

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
)

type (
	Expression struct {
		indexName *string
		hashKey   *string
		rangeKey  *string

		item interface{}

		expressions map[string]domain.WithCondition
	}
)

func NewSqlBuilder(config *domain.Config) domain.SqlExpression {
	return &Expression{
		hashKey:     aws.String(config.Table.GetMetadata().GetHash()),
		rangeKey:    aws.String(config.Table.GetMetadata().GetRange()),
		expressions: map[string]domain.WithCondition{},
	}
}

/* Expression */

func (e *Expression) SetIndex(indexName string) domain.SqlExpression {
	e.indexName = aws.String(indexName)
	return e
}

func (e *Expression) IndexName() *string {
	return e.indexName
}

func (e *Expression) Where(condition domain.WithCondition) domain.SqlExpression {
	e.expressions["key"] = condition
	return e
}

func (e *Expression) AndWhere(keyCondition domain.WithSortKeyCondition) domain.SqlExpression {
	if keyCondition.HasSortKey() {
		e.expressions["sortKey"] = keyCondition
	}

	return e
}

func (e *Expression) ExpressionAttributeValues() map[string]types.AttributeValue {
	if e.expressions["sortKey"] == nil {
		return map[string]types.AttributeValue{":key": e.expressions["key"].Value()}
	}

	sortKeyCondition := e.expressions["sortKey"].(domain.WithSortKeyCondition)

	if sortKeyCondition.SimpleCondition() {
		return map[string]types.AttributeValue{
			":key":     e.expressions["key"].Value(),
			":sortVal": sortKeyCondition.Value(),
		}
	}

	return map[string]types.AttributeValue{
		":key": e.expressions["key"].Value(),
		":start": e.getAttributeValueMember(
			reflect.ValueOf(sortKeyCondition.StartValue()),
		),
		":end": e.getAttributeValueMember(
			reflect.ValueOf(sortKeyCondition.EndValue()),
		),
	}
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
	if e.expressions["sortKey"] == nil {
		return aws.String(e.expressions["key"].KeyCondition())
	}

	sortKeyCondition := e.expressions["sortKey"].KeyCondition()

	return aws.String(fmt.Sprintf(
		"%s and %s",
		e.expressions["key"].KeyCondition(),
		sortKeyCondition,
	))
}

func (e *Expression) SetItem(item interface{}) domain.SqlExpression {
	if reflect.TypeOf(item).Kind() != reflect.Struct {
		panic("item should be a struct")
	}

	e.item = item

	return e
}

func (e *Expression) Names() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{}
}

func (e *Expression) Values() map[string]types.AttributeValue {
	if e.item == nil {
		panic("to call Values, before call SetItem")
	}

	item := reflect.ValueOf(e.item)

	attributes := map[string]types.AttributeValue{}

	for i := 0; i < item.NumField(); i++ {
		field := item.Field(i)
		name := item.Type().Field(i).Name

		attributes[name] = e.getAttributeValueMember(field)
	}

	return attributes
}

func (e *Expression) getAttributeValueMember(val reflect.Value) types.AttributeValue {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &types.AttributeValueMemberN{Value: val.String()}
	case reflect.Bool:
		return &types.AttributeValueMemberBOOL{Value: val.Bool()}
	//case reflect.Map:

	case reflect.Slice, reflect.Array:
		return &types.AttributeValueMemberSS{Value: val.Interface().([]string)}
	default:
		return &types.AttributeValueMemberS{Value: val.String()}
	}
}
