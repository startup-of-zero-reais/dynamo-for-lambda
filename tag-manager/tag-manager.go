package tagManager

import (
	"github.com/startup-of-zero-reais/dynamo-for-lambda/logger"
	"reflect"
)

type (
	// Manager é o contrato de um gerenciador do pacote tagManager
	Manager interface {
		SetEntity(entity interface{}) Manager
		MapTags() error
		GetMapper() TagMapperInterface

		TagGetters
	}

	// TagManager é a estrutura que gerencia as tags extraídas.
	// Tags específicas da chave: diinamo
	TagManager struct {
		StructToMap interface{}
		TagMapper   TagMapperInterface
		logger.Log
	}

	TagGetters interface {
		GetHash() string
		GetRange() string
		GetType(key string) reflect.Kind
	}
)

// NewTagManager inicializa uma estrutura TagManager
func NewTagManager() *TagManager {
	logg := logger.NewLogger()

	return &TagManager{
		Log: logg,
		TagMapper: &TagMapper{
			Log: logg,
		},
	}
}

// SetEntity define uma estrutura que será iterada.
// A estrutura não precisa ter valores preenchidos, deve apenas, ter as
// tags que correspondem a chave diinamo
func (t *TagManager) SetEntity(entity interface{}) Manager {
	t.StructToMap = entity

	// Inicializa o TagMapper dentro de TagManager
	t.TagMapper.SetPropertyTypes(reflect.TypeOf(t.StructToMap))

	return t
}

// MapTags é um método que executa a extração das tags e faz o
// mapeamento do TagMapper para o TagManager
func (t *TagManager) MapTags() error {
	return t.TagMapper.RunMap()
}

// GetHash devolve o nome do campo que representa o Hash da tabela
func (t *TagManager) GetHash() string {
	return t.TagMapper.GetHash()
}

// GetRange devolve o nome do campo que representa a Sort Key da tabela
func (t *TagManager) GetRange() string {
	return t.TagMapper.GetRange()
}

// GetType recupera o valor definido pela tag type de uma key específica
func (t *TagManager) GetType(key string) reflect.Kind {
	return t.TagMapper.GetType(key)
}

// GetMapper é um método para recuperar o TagMapper
// dentro de TagManager
func (t *TagManager) GetMapper() TagMapperInterface {
	return t.TagMapper
}

/* Exemplo de entidade e tags aceitas */

// ExampleEntity é um exemplo de entidade com as tags aceitas
// e exemplos de como definir as tags.
//
// 	OBS: ExampleEntity não deve ser usado!
type ExampleEntity struct {
	PK           int    `diinamo:"type:number;hash"`
	SK           string `diinamo:"type:string;range"`
	Owner        string `diinamo:"type:string;gsi:CourseOwnerIndex;keyPairs:PK=Owner"`
	Title        string `diinamo:"type:string;gsi:CourseTitleIndex;keyPairs:Title=SK"`
	ParentCourse string `diinamo:"type:string;gsi:CourseLessonsIndex;keyPairs:ParentCourse=SK"`
	ParentModule string `diinamo:"type:string;lsi:ModuleLessonsIndex;keyPairs:ParentModule=SK"`
}

// This é um método apenas de teste para ExampleEntity
func (e *ExampleEntity) This() *ExampleEntity {
	return e
}
