package domain

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type (
	Action      string
	Environment string

	Dynamo interface {
		Perform(action Action, sql SqlExpression, result interface{}) error
		NewExpressionBuilder() SqlExpression
		Migrate() error
		Seed(items ...*dynamodb.PutItemInput) error
	}
)

func (e Environment) IsDev() bool {
	return string(e) == "development"
}
