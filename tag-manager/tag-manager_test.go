package tagManager_test

import (
	mocks "github.com/startup-of-zero-reais/dynamo-for-lambda/mocks/tag-manager"
	tagManager "github.com/startup-of-zero-reais/dynamo-for-lambda/tag-manager"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/tag-manager/logger"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestTagManager_SetEntity(t *testing.T) {
	t.Run("should update StructToMap entity", func(t *testing.T) {
		logg := logger.NewLogger()
		tm := &tagManager.TagManager{
			Log:       logg,
			TagMapper: &tagManager.TagMapper{},
		}
		assert.Nil(t, tm.StructToMap)

		tm.SetEntity(mocks.Mocktable{})
		assert.NotNil(t, tm.StructToMap)
	})
}

func TestTagManager_MapTags(t *testing.T) {
	t.Run("should call RunMap inside TagMapper", func(t *testing.T) {
		logg := logger.NewLogger()
		tagMapper := new(mocks.TagMapperInterface)
		tagMapper.On("SetPropertyTypes", reflect.TypeOf(mocks.Mocktable{})).Return()
		tagMapper.On("RunMap").Return(nil)

		tm := &tagManager.TagManager{
			Log:       logg,
			TagMapper: tagMapper,
		}

		tm.SetEntity(mocks.Mocktable{})

		err := tm.MapTags()

		assert.Nil(t, err)
	})
}

func TestTagManager_TagGetters(t *testing.T) {
	t.Run("should call TagGetters", func(t *testing.T) {
		tagMapper := new(mocks.TagMapperInterface)
		tagMapper.On("GetHash").Return("PK")
		tagMapper.On("GetRange").Return("SK")
		tagMapper.On("GetType", "PK").Return(reflect.String)

		tm := tagManager.NewTagManager()
		tm.TagMapper = tagMapper

		hash := tm.GetHash()
		_range := tm.GetRange()
		_type := tm.GetType("PK")

		assert.NotZero(t, hash)
		assert.NotZero(t, _range)
		assert.NotZero(t, _type)
	})
}
