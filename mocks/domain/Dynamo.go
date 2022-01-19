// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	domain "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"

	mock "github.com/stretchr/testify/mock"
)

// Dynamo is an autogenerated mock type for the Dynamo type
type Dynamo struct {
	mock.Mock
}

// Migrate provides a mock function with given fields:
func (_m *Dynamo) Migrate() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewExpressionBuilder provides a mock function with given fields:
func (_m *Dynamo) NewExpressionBuilder() domain.SqlExpression {
	ret := _m.Called()

	var r0 domain.SqlExpression
	if rf, ok := ret.Get(0).(func() domain.SqlExpression); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.SqlExpression)
		}
	}

	return r0
}

// Perform provides a mock function with given fields: action, sql, result
func (_m *Dynamo) Perform(action domain.Action, sql domain.SqlExpression, result interface{}) error {
	ret := _m.Called(action, sql, result)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.Action, domain.SqlExpression, interface{}) error); ok {
		r0 = rf(action, sql, result)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Seed provides a mock function with given fields: items
func (_m *Dynamo) Seed(items ...*dynamodb.PutItemInput) error {
	_va := make([]interface{}, len(items))
	for _i := range items {
		_va[_i] = items[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...*dynamodb.PutItemInput) error); ok {
		r0 = rf(items...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}