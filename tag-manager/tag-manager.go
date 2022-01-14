package tagManager

import (
	"github.com/startup-of-zero-reais/dynamo-for-lambda/tag-manager/logger"
	"reflect"
)

type (
	// Manager é o contrato de um gerenciador do pacote tagManager
	Manager interface {
		SetEntity(entity interface{}) Manager
		MapTags() error
	}

	// TagManager é a estrutura que gerencia as tags extraídas.
	// Tags específicas da chave: diinamo
	TagManager struct {
		StructToMap interface{}
		Tags        map[string]interface{}

		logger.Log

		*TagMapper
	}
)

// NewTagManager inicializa uma estrutura TagManager
func NewTagManager() *TagManager {
	logger := logger.NewLogger()

	return &TagManager{
		Log: logger,
		TagMapper: &TagMapper{
			Log: logger,
		},
	}
}

// SetEntity define uma estrutura que será iterada.
// A estrutura não precisa ter valores preenchidos, deve apenas, ter as
// tags que correspondem a chave diinamo
func (t *TagManager) SetEntity(entity interface{}) *TagManager {
	t.StructToMap = entity

	// Inicializa o TagMapper dentro de TagManager
	t.TagMapper.PropertyTypes = reflect.TypeOf(t.StructToMap)

	return t
}

// MapTags é um método que executa a extração das tags e faz o
// mapeamento do TagMapper para o TagManager
func (t *TagManager) MapTags() error {
	return t.RunMap()
}
