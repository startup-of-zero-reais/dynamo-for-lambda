package domain

type (
	DynamoSQL interface {
		Get(expression SqlExpression, target interface{}) error
		Put(item interface{}, result interface{}) error
		Update(expression interface{}, item interface{}, result interface{}) error
		Delete(expression SqlExpression) error
	}
)
