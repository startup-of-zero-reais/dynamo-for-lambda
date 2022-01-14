package tagManager_test

import (
	"errors"
	"fmt"
	mocks "github.com/startup-of-zero-reais/dynamo-for-lambda/mocks/tag-manager"
	tagManager "github.com/startup-of-zero-reais/dynamo-for-lambda/tag-manager"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/tag-manager/logger"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func prepareTagMapper() *tagManager.TagMapper {
	tm := &tagManager.TagMapper{
		PropertyTypes: reflect.TypeOf(mocks.Mocktable{}),
		Log:           logger.NewLogger(),
		TagsModel:     new(tagManager.TagsModel),
	}

	tm.ExtractFieldList()

	return tm
}

func TestTagMapper_ExtractFieldList(t *testing.T) {
	t.Run("should extract field list from PropertyTypes", func(t *testing.T) {
		tagMapper := &tagManager.TagMapper{
			PropertyTypes: reflect.TypeOf(mocks.Mocktable{}),
		}

		assert.Nil(t, tagMapper.FieldList)
		assert.Nil(t, tagMapper.FieldNames)

		tagMapper.ExtractFieldList()

		assert.NotNil(t, tagMapper.FieldList)
		assert.NotNil(t, tagMapper.FieldNames)
	})
}

func TestTagMapper_RunMap(t *testing.T) {
	t.Run("should run mapper", func(t *testing.T) {
		tm := &tagManager.TagMapper{
			PropertyTypes: reflect.TypeOf(mocks.Mocktable{}),
			Log:           logger.NewLogger(),
		}

		err := tm.RunMap()
		assert.Nil(t, err)
		assert.NotNil(t, tm.PropertyTypes)
		assert.NotNil(t, tm.FieldNames)
		assert.NotNil(t, tm.FieldList)
		assert.NotNil(t, tm.TagsModel)
		assert.NotNil(t, tm.Hash)
		assert.NotNil(t, tm.Range)
		assert.NotNil(t, tm.GSI)
		assert.NotNil(t, tm.LSI)
		assert.NotNil(t, tm.Types)
	})
	t.Run("should capture fail tags loop", func(t *testing.T) {
		tm := &tagManager.TagMapper{
			PropertyTypes: reflect.TypeOf(mocks.Mocktable{}),
			Log:           logger.NewLogger(),
		}

		tm.TagsModel = &tagManager.TagsModel{}

		for i := 0; i < 5; i++ {
			tm.LSI = append(tm.LSI, tagManager.LocalSecIndex{
				IndexName: fmt.Sprintf("index-%d", i),
				Hash:      fmt.Sprintf("hash-%d", i),
				Range:     fmt.Sprintf("range-%d", i),
			})
		}

		err := tm.RunMap()
		assert.EqualError(t, err, "max local secondary index reached")
	})
}

func TestTagMapper_TagsLoop(t *testing.T) {
	t.Run("should iterate on field list", func(t *testing.T) {
		tm := prepareTagMapper()
		err := tm.TagsLoop(tm.ExtractPK)

		assert.Nil(t, err)
	})
	t.Run("should fail with an TagHandler", func(t *testing.T) {
		tm := prepareTagMapper()

		if len(tm.FieldList) <= 0 {
			t.Fatalf("tagMapper has an empty field list")
		}

		th := new(mocks.TagHandler)
		tagPairs := []string{"type:string", "hash"}
		th.On("Execute", tagPairs, tm.FieldList[0]).Return(errors.New("fail tag handler"))

		err := tm.TagsLoop(th.Execute)

		assert.EqualError(t, err, "fail tag handler")
	})
}

func TestTagMapper_ExtractPK(t *testing.T) {
	t.Run("should extract PK", func(t *testing.T) {
		tm := prepareTagMapper()
		tm.TagsModel = nil

		for _, f := range tm.FieldList {
			if inlineTags, ok := f.Tag.Lookup("diinamo"); ok {
				tags := strings.Split(inlineTags, ";")
				err := tm.ExtractPK(tags, f)
				assert.Nil(t, err)
			}
		}

		assert.NotNil(t, tm.Hash)
		assert.NotNil(t, tm.Range)
	})
}

func TestTagMapper_ExtractGSI(t *testing.T) {
	t.Run("should extract gsi", func(t *testing.T) {
		tm := prepareTagMapper()

		for _, f := range tm.FieldList {
			if inlineTags, ok := f.Tag.Lookup("diinamo"); ok {
				tags := strings.Split(inlineTags, ";")
				err := tm.ExtractGSI(tags, f)
				assert.Nil(t, err)
			}
		}

		assert.NotNil(t, tm.TagsModel.GSI)
		assert.True(t, len(tm.TagsModel.GSI) > 0)
	})
	t.Run("should fail if has more than 20 gsi", func(t *testing.T) {
		tm := prepareTagMapper()

		for i := 0; i < 20; i++ {
			tm.GSI = append(tm.GSI, tagManager.GlobalSecIndex{
				IndexName: fmt.Sprintf("index-%d", i),
				Hash:      fmt.Sprintf("hash-%d", i),
				Range:     fmt.Sprintf("range-%d", i),
			})
		}

		for _, f := range tm.FieldList {
			if inlineTags, ok := f.Tag.Lookup("diinamo"); ok {
				tags := strings.Split(inlineTags, ";")
				err := tm.ExtractGSI(tags, f)
				assert.EqualError(t, err, "max global secondary index reached")
			}
		}
	})
}

func TestTagMapper_ExtractLSI(t *testing.T) {
	t.Run("should extract lsi", func(t *testing.T) {
		tm := prepareTagMapper()

		for _, f := range tm.FieldList {
			if inlineTags, ok := f.Tag.Lookup("diinamo"); ok {
				tags := strings.Split(inlineTags, ";")
				err := tm.ExtractLSI(tags, f)
				assert.Nil(t, err)
			}
		}

		assert.NotNil(t, tm.TagsModel.LSI)
		assert.True(t, len(tm.TagsModel.LSI) > 0)
	})
	t.Run("should fail if has more than 5 lsi", func(t *testing.T) {
		tm := prepareTagMapper()

		for i := 0; i < 5; i++ {
			tm.LSI = append(tm.LSI, tagManager.LocalSecIndex{
				IndexName: fmt.Sprintf("index-%d", i),
				Hash:      fmt.Sprintf("hash-%d", i),
				Range:     fmt.Sprintf("range-%d", i),
			})
		}

		for _, f := range tm.FieldList {
			if inlineTags, ok := f.Tag.Lookup("diinamo"); ok {
				tags := strings.Split(inlineTags, ";")
				err := tm.ExtractLSI(tags, f)
				assert.EqualError(t, err, "max local secondary index reached")
			}
		}
	})
}

func TestTagMapper_ExtractTypes(t *testing.T) {
	t.Run("should extract types tag definition", func(t *testing.T) {
		tm := prepareTagMapper()

		var tagsPair []string
		var field reflect.StructField
		for _, f := range tm.FieldList {
			if inlineTags, ok := f.Tag.Lookup("diinamo"); ok {
				tags := strings.Split(inlineTags, ";")
				tagsPair = tags
				field = f
			}
		}

		err := tm.ExtractTypes(tagsPair, field)
		assert.Nil(t, err)
		assert.NotNil(t, tm.Types)
	})
}
