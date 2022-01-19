package tagManager_test

import (
	"fmt"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/logger"
	tableMock "github.com/startup-of-zero-reais/dynamo-for-lambda/mocks/domain"
	mocks "github.com/startup-of-zero-reais/dynamo-for-lambda/mocks/tag-manager"
	tagManager "github.com/startup-of-zero-reais/dynamo-for-lambda/tag-manager"
	"github.com/stretchr/testify/assert"
	"log"
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

		tm.SetEntity(tableMock.Mocktable{})
		assert.NotNil(t, tm.StructToMap)
	})
}

func TestTagManager_MapTags(t *testing.T) {
	t.Run("should call RunMap inside TagMapper", func(t *testing.T) {
		logg := logger.NewLogger()
		tagMapper := new(mocks.TagMapperInterface)
		tagMapper.On("SetPropertyTypes", reflect.TypeOf(tableMock.Mocktable{})).Return()
		tagMapper.On("RunMap").Return(nil)

		tm := &tagManager.TagManager{
			Log:       logg,
			TagMapper: tagMapper,
		}

		tm.SetEntity(tableMock.Mocktable{})

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

func TestExampleEntity_This(t *testing.T) {
	t.Run("should return this", func(t *testing.T) {
		ee := &tagManager.ExampleEntity{}
		e := ee.This()

		assert.Equal(t, ee, e)
	})
}

func ExampleTagManager_GetHash() {
	// Nova instância de TagManager
	tm := tagManager.NewTagManager()

	// Define a entidade a ter as tags mapeadas e imediatamente
	// mapeia as tags da entidade
	// MapTags retorna um erro
	err := tm.SetEntity(tableMock.Mocktable{}).MapTags()
	if err != nil {
		// Tratamento de erro
		log.Fatalln(err)
	}

	fmt.Println(tm.GetHash())
	// Output:
	// PK
}

func ExampleTagManager_GetRange() {
	// Nova instância de TagManager
	tm := tagManager.NewTagManager()

	// Define a entidade a ter as tags mapeadas e imediatamente
	// mapeia as tags da entidade
	// MapTags retorna um erro
	err := tm.SetEntity(tableMock.Mocktable{}).MapTags()
	if err != nil {
		// Tratamento de erro
		log.Fatalln(err)
	}

	fmt.Println(tm.GetRange())
	// Output:
	// SK
}

func ExampleTagManager_GetType() {
	// Nova instância de TagManager
	tm := tagManager.NewTagManager()

	// Define a entidade a ter as tags mapeadas e imediatamente
	// mapeia as tags da entidade
	// MapTags retorna um erro
	err := tm.SetEntity(tableMock.Mocktable{}).MapTags()
	if err != nil {
		// Tratamento de erro
		log.Fatalln(err)
	}

	// GetType retorna um valor do tipo reflect.Kind
	fmt.Println(tm.GetType("PK"))
	// Output:
	// string
}

func ExampleTagManager_MapTags() {
	// Nova instância de TagManager
	tm := tagManager.NewTagManager()

	// Define a entidade a ter as tags mapeadas e imediatamente
	// mapeia as tags da entidade
	// MapTags retorna um erro
	err := tm.SetEntity(tableMock.Mocktable{}).MapTags()
	if err != nil {
		// Tratamento de erro
		log.Fatalln(err)
	}

	fmt.Println(tm.TagMapper.GetHash())
	// Output:
	// PK
}

func ExampleTagManager_SetEntity() {
	// Nova instância de TagManager
	tm := tagManager.NewTagManager()

	// Define a entidade a ter as tags mapeadas
	tm.SetEntity(tagManager.ExampleEntity{})
	// Output:
}
