// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	reflect "reflect"

	mock "github.com/stretchr/testify/mock"

	tagManager "github.com/startup-of-zero-reais/dynamo-for-lambda/tag-manager"
)

// TagMapperInterface is an autogenerated mock type for the TagMapperInterface type
type TagMapperInterface struct {
	mock.Mock
}

// ExtractFieldList provides a mock function with given fields:
func (_m *TagMapperInterface) ExtractFieldList() {
	_m.Called()
}

// ExtractGSI provides a mock function with given fields: tagsPair, field
func (_m *TagMapperInterface) ExtractGSI(tagsPair []string, field reflect.StructField) error {
	ret := _m.Called(tagsPair, field)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, reflect.StructField) error); ok {
		r0 = rf(tagsPair, field)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExtractLSI provides a mock function with given fields: tagsPair, field
func (_m *TagMapperInterface) ExtractLSI(tagsPair []string, field reflect.StructField) error {
	ret := _m.Called(tagsPair, field)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, reflect.StructField) error); ok {
		r0 = rf(tagsPair, field)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExtractPK provides a mock function with given fields: tagsPair, field
func (_m *TagMapperInterface) ExtractPK(tagsPair []string, field reflect.StructField) error {
	ret := _m.Called(tagsPair, field)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, reflect.StructField) error); ok {
		r0 = rf(tagsPair, field)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExtractTypes provides a mock function with given fields: tagsPair, field
func (_m *TagMapperInterface) ExtractTypes(tagsPair []string, field reflect.StructField) error {
	ret := _m.Called(tagsPair, field)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, reflect.StructField) error); ok {
		r0 = rf(tagsPair, field)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetHash provides a mock function with given fields:
func (_m *TagMapperInterface) GetHash() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetModel provides a mock function with given fields:
func (_m *TagMapperInterface) GetModel() *tagManager.TagsModel {
	ret := _m.Called()

	var r0 *tagManager.TagsModel
	if rf, ok := ret.Get(0).(func() *tagManager.TagsModel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tagManager.TagsModel)
		}
	}

	return r0
}

// GetRange provides a mock function with given fields:
func (_m *TagMapperInterface) GetRange() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetType provides a mock function with given fields: key
func (_m *TagMapperInterface) GetType(key string) reflect.Kind {
	ret := _m.Called(key)

	var r0 reflect.Kind
	if rf, ok := ret.Get(0).(func(string) reflect.Kind); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(reflect.Kind)
	}

	return r0
}

// RunMap provides a mock function with given fields:
func (_m *TagMapperInterface) RunMap() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetPropertyTypes provides a mock function with given fields: v
func (_m *TagMapperInterface) SetPropertyTypes(v reflect.Type) {
	_m.Called(v)
}

// TagsLoop provides a mock function with given fields: cases
func (_m *TagMapperInterface) TagsLoop(cases ...tagManager.TagHandler) error {
	_va := make([]interface{}, len(cases))
	for _i := range cases {
		_va[_i] = cases[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...tagManager.TagHandler) error); ok {
		r0 = rf(cases...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
