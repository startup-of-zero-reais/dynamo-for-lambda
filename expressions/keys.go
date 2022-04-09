package expressions

import (
	"fmt"
	"log"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
)

type (
	condition struct {
		expression string
		condition  domain.Condition
	}

	KeyCondition struct {
		name *string
		Val  interface{}
	}

	SortKeyCondition struct {
		exists bool
		name   *string
		Val    interface{}

		betweenStart interface{}
		betweenEnd   interface{}

		condition condition
	}
)

/* KeyCondition */

func NewKeyCondition(name string, val interface{}) *KeyCondition {
	kc := &KeyCondition{Val: val}
	kc.SetName(name)

	return kc
}

func (k *KeyCondition) SetName(name string) domain.WithCondition {
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

	return GetAttributeValueMemberType(reflect.ValueOf(k.Val))
}

func (k *KeyCondition) KeyCondition() string {
	return fmt.Sprintf("%s = :key", *k.name)
}

/* SortKeyCondition */

func NewSortKeyCondition(name string) domain.WithSortKeyCondition {
	ksc := &SortKeyCondition{}
	ksc.SetName(name)

	return ksc
}

func (k *SortKeyCondition) SetName(name string) domain.WithCondition {
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

func (k *SortKeyCondition) HasSortKey() bool {
	return k.exists
}

func (k *SortKeyCondition) StarsWith(value interface{}) domain.WithSortKeyCondition {
	k.condition = condition{
		expression: fmt.Sprintf("begins_with(%s, :sortVal)", *k.name),
		condition:  StartsWith,
	}
	k.Val = value
	return k
}

func (k *SortKeyCondition) Equal(value interface{}) domain.WithSortKeyCondition {
	k.condition = condition{
		expression: fmt.Sprintf("%s = :sortVal", *k.name),
		condition:  Equal,
	}
	k.Val = value
	return k
}

func (k *SortKeyCondition) LessThan(value interface{}) domain.WithSortKeyCondition {
	k.condition = condition{
		expression: fmt.Sprintf("%s < :sortVal", *k.name),
		condition:  LessThan,
	}
	k.Val = value
	return k
}

func (k *SortKeyCondition) LessThanOrEqual(value interface{}) domain.WithSortKeyCondition {
	k.condition = condition{
		expression: fmt.Sprintf("%s <= :sortVal", *k.name),
		condition:  LessThanOrEqual,
	}
	k.Val = value

	return k
}

func (k *SortKeyCondition) GreaterThan(value interface{}) domain.WithSortKeyCondition {
	k.condition = condition{
		expression: fmt.Sprintf("%s > :sortVal", *k.name),
		condition:  GreaterThan,
	}
	k.Val = value

	return k
}

func (k *SortKeyCondition) GreaterThanOrEqual(value interface{}) domain.WithSortKeyCondition {
	k.condition = condition{
		expression: fmt.Sprintf("%s >= :sortVal", *k.name),
		condition:  GreaterThanOrEqual,
	}
	k.Val = value

	return k
}

func (k *SortKeyCondition) Between(start, end interface{}) domain.WithSortKeyCondition {
	k.condition = condition{
		expression: fmt.Sprintf("%sBETWEEN:startAND:end", *k.name),
		condition:  Between,
	}
	k.betweenStart = start
	k.betweenEnd = end

	return k
}

func (k *SortKeyCondition) KeyCondition() string {
	return k.condition.expression
}

func (k *SortKeyCondition) SimpleCondition() bool {
	return k.condition.condition != Between
}

func (k *SortKeyCondition) StartValue() interface{} {
	return k.betweenStart
}

func (k *SortKeyCondition) EndValue() interface{} {
	return k.betweenEnd
}
