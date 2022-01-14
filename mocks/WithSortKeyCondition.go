// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	codDynamo "github.com/startup-of-zero-reais/dynamo-for-lambda"
	mock "github.com/stretchr/testify/mock"

	types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// WithSortKeyCondition is an autogenerated mock type for the WithSortKeyCondition type
type WithSortKeyCondition struct {
	mock.Mock
}

// Name provides a mock function with given fields:
func (_m *WithSortKeyCondition) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SetName provides a mock function with given fields: name
func (_m *WithSortKeyCondition) SetName(name string) codDynamo.WithCondition {
	ret := _m.Called(name)

	var r0 codDynamo.WithCondition
	if rf, ok := ret.Get(0).(func(string) codDynamo.WithCondition); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(codDynamo.WithCondition)
		}
	}

	return r0
}

// Value provides a mock function with given fields:
func (_m *WithSortKeyCondition) Value() types.AttributeValue {
	ret := _m.Called()

	var r0 types.AttributeValue
	if rf, ok := ret.Get(0).(func() types.AttributeValue); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.AttributeValue)
		}
	}

	return r0
}

// hasSortKey provides a mock function with given fields:
func (_m *WithSortKeyCondition) hasSortKey() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}