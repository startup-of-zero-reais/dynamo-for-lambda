package expressions

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
)

type (
	KeyCondition struct {
		name *string
		Val  interface{}
	}

	SortKeyCondition struct {
		exists bool
		name   *string
		Val    interface{}
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

	return &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", k.Val)}
}

/* SortKeyCondition */

func NewSortKeyCondition(name string, val interface{}) domain.WithSortKeyCondition {
	ksc := &SortKeyCondition{Val: val}
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
	return k.exists == true
}
