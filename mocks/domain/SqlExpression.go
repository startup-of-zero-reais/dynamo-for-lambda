// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	domain "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	mock "github.com/stretchr/testify/mock"

	types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// SqlExpression is an autogenerated mock type for the SqlExpression type
type SqlExpression struct {
	mock.Mock
}

// AndWhere provides a mock function with given fields: keyCondition
func (_m *SqlExpression) AndWhere(keyCondition domain.WithSortKeyCondition) domain.SqlExpression {
	ret := _m.Called(keyCondition)

	var r0 domain.SqlExpression
	if rf, ok := ret.Get(0).(func(domain.WithSortKeyCondition) domain.SqlExpression); ok {
		r0 = rf(keyCondition)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.SqlExpression)
		}
	}

	return r0
}

// Build provides a mock function with given fields:
func (_m *SqlExpression) Build() map[string]types.AttributeValue {
	ret := _m.Called()

	var r0 map[string]types.AttributeValue
	if rf, ok := ret.Get(0).(func() map[string]types.AttributeValue); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]types.AttributeValue)
		}
	}

	return r0
}

// Key provides a mock function with given fields:
func (_m *SqlExpression) Key() map[string]types.AttributeValue {
	ret := _m.Called()

	var r0 map[string]types.AttributeValue
	if rf, ok := ret.Get(0).(func() map[string]types.AttributeValue); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]types.AttributeValue)
		}
	}

	return r0
}

// KeyCondition provides a mock function with given fields:
func (_m *SqlExpression) KeyCondition() *string {
	ret := _m.Called()

	var r0 *string
	if rf, ok := ret.Get(0).(func() *string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	return r0
}

// Names provides a mock function with given fields:
func (_m *SqlExpression) Names() map[string]types.AttributeValue {
	ret := _m.Called()

	var r0 map[string]types.AttributeValue
	if rf, ok := ret.Get(0).(func() map[string]types.AttributeValue); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]types.AttributeValue)
		}
	}

	return r0
}

// SetIndex provides a mock function with given fields: indexName
func (_m *SqlExpression) SetIndex(indexName string) domain.SqlExpression {
	ret := _m.Called(indexName)

	var r0 domain.SqlExpression
	if rf, ok := ret.Get(0).(func(string) domain.SqlExpression); ok {
		r0 = rf(indexName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.SqlExpression)
		}
	}

	return r0
}

// Values provides a mock function with given fields:
func (_m *SqlExpression) Values() map[string]types.AttributeValue {
	ret := _m.Called()

	var r0 map[string]types.AttributeValue
	if rf, ok := ret.Get(0).(func() map[string]types.AttributeValue); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]types.AttributeValue)
		}
	}

	return r0
}

// Where provides a mock function with given fields: condition
func (_m *SqlExpression) Where(condition domain.WithCondition) domain.SqlExpression {
	ret := _m.Called(condition)

	var r0 domain.SqlExpression
	if rf, ok := ret.Get(0).(func(domain.WithCondition) domain.SqlExpression); ok {
		r0 = rf(condition)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.SqlExpression)
		}
	}

	return r0
}
