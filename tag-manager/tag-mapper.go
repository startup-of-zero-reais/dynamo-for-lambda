package tagManager

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/startup-of-zero-reais/dynamo-for-lambda/logger"
)

type (
	// TagHandler é uma assinatura de um handler de tags.
	// Usado para manipular as diferentes tags recebidas pela chave diinamo
	TagHandler func(tagsPair []string, field reflect.StructField) error

	// ProvisionedThroughput é uma estrutura de definição de Throughput de
	// provisionamento do DynamoDB
	ProvisionedThroughput struct {
		ReadCapacity  uint
		WriteCapacity uint
	}

	// GlobalSecIndex é o formato de um Global Secondary Index (GSI).
	// Essa estrutura mantém o IndexName, e o par de chaves Hash e Range
	GlobalSecIndex struct {
		IndexName string
		Hash      string
		Range     string
		ProvisionedThroughput
	}

	// LocalSecIndex é o formato de um Local Secondary Index (LSI).
	// Essa estrutura mantém o IndexName, e o par de chaves Hash e Range
	LocalSecIndex struct {
		IndexName string
		Hash      string
		Range     string
		ProvisionedThroughput
	}

	// TagsModel é uma estrutura de gerenciamento de tags.
	// Responsável por manter as Hash, RangeKeys, GSI, Types, etc.
	TagsModel struct {
		Hash  string
		Range string

		GSI []GlobalSecIndex
		LSI []LocalSecIndex

		Types map[string]reflect.Kind
	}

	// TagMapper é uma estrutura para gerenciar os dados das tags
	TagMapper struct {
		PropertyTypes reflect.Type

		FieldNames []string
		FieldList  []reflect.StructField

		*TagsModel
		logger.Log
	}

	// TagMapperInterface sugere o contrato de TagMapper com
	// o fim de manter a consistência de métodos e retornos
	TagMapperInterface interface {
		SetPropertyTypes(v reflect.Type)
		ExtractFieldList()
		RunMap() error
		TagsLoop(cases ...TagHandler) error
		ExtractPK(tagsPair []string, field reflect.StructField) error
		ExtractGSI(tagsPair []string, field reflect.StructField) error
		ExtractLSI(tagsPair []string, field reflect.StructField) error
		ExtractTypes(tagsPair []string, field reflect.StructField) error

		GetModel() *TagsModel

		TagGetters
	}
)

// Tags
const (
	hash     = "hash"
	_range   = "range"
	gsi      = "gsi"
	keyPairs = "keyPairs"
	lsi      = "lsi"
	_type    = "type"
)

// ExtractFieldList extrai os metadados de PropertyTypes de TagMapper
// para a estrutura.
func (t *TagMapper) ExtractFieldList() {
	for i := 0; i < t.PropertyTypes.NumField(); i++ {
		field := t.PropertyTypes.Field(i)
		t.FieldNames = append(t.FieldNames, field.Name)
		t.FieldList = append(t.FieldList, field)
	}
}

// RunMap é o método que irá criar, manipular e definir as tags
// em TagsModel dentro de TagMapper
func (t *TagMapper) RunMap() error {
	started := time.Now()
	t.Debug("running map...")
	t.ExtractFieldList()

	if err := t.TagsLoop(
		t.ExtractPK,
		t.ExtractGSI,
		t.ExtractLSI,
		t.ExtractTypes,
	); err != nil {
		return err
	}

	t.Debug("extraction complete. %v spent\n", time.Since(started))

	return nil
}

// TagsLoop é o método que faz a iteração nos campos da Struct recebida
// em TagManager.SetEntity do TagManager
func (t *TagMapper) TagsLoop(cases ...TagHandler) error {
	for _, field := range t.FieldList {
		if inlineTags, ok := field.Tag.Lookup("diinamo"); ok {
			// inlineTags se parece com:
			// 	type:string;hash
			tagsPair := strings.Split(inlineTags, ";")

			if len(cases) > 0 {
				for _, useCase := range cases {
					err := useCase(tagsPair, field)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// ExtractPK é um método de extração e definição dos pares de Chave:
// Hash e Range
func (t *TagMapper) ExtractPK(tagsPair []string, field reflect.StructField) error {
	for _, tag := range tagsPair {
		if t.TagsModel == nil {
			t.TagsModel = &TagsModel{}
		}

		switch tag {
		case hash:
			t.TagsModel.Hash = field.Name
		case _range:
			t.TagsModel.Range = field.Name
		}
	}

	return nil
}

// ExtractGSI é um método de extração e definição dos pares de Chave
// dos índices secundários: Hash, Range e IndexName
func (t *TagMapper) ExtractGSI(tagsPair []string, field reflect.StructField) error {
	gsIndex := &GlobalSecIndex{
		ProvisionedThroughput: ProvisionedThroughput{
			ReadCapacity:  1,
			WriteCapacity: 1,
		},
	}

	// tagsPair se parece com:
	// 	[gsi:TheNameOfIndex keyPairs:PK=SK]
	// 	[gsi:TheNameOfAnotherIndex keyPairs:SK=PK]
	for _, tag := range tagsPair {
		tagKeyValue := strings.Split(tag, ":")

		if t.TagsModel.GSI == nil {
			t.TagsModel.GSI = []GlobalSecIndex{}
		}

		switch tagKeyValue[0] {
		case gsi:
			gsIndex.IndexName = tagKeyValue[1]
		case keyPairs:
			hashRange := strings.Split(tagKeyValue[1], "=")
			gsIndex.Hash = hashRange[0]
			gsIndex.Range = hashRange[1]
		}
	}

	if len(t.TagsModel.GSI) >= 20 {
		return errors.New("max global secondary index reached")
	}

	t.Info("adding gsi: %+v\n", gsIndex)
	if gsIndex.IndexName != "" && gsIndex.Hash != "" {
		t.TagsModel.GSI = append(t.TagsModel.GSI, *gsIndex)
	}

	return nil
}

// ExtractLSI é um método de extração de definição dos pares de Chave
// dos índices secundários locais: Hash, Range, IndexName e Projection
func (t *TagMapper) ExtractLSI(tagsPair []string, field reflect.StructField) error {
	lsIndex := &LocalSecIndex{
		ProvisionedThroughput: ProvisionedThroughput{ReadCapacity: 1, WriteCapacity: 1},
	}

	for _, tag := range tagsPair {
		tagKeyValue := strings.Split(tag, ":")

		if t.TagsModel.LSI == nil {
			t.TagsModel.LSI = []LocalSecIndex{}
		}

		switch tagKeyValue[0] {
		case lsi:
			lsIndex.IndexName = tagKeyValue[1]
		case keyPairs:
			hashRange := strings.Split(tagKeyValue[1], "=")
			lsIndex.Hash = hashRange[0]
			lsIndex.Range = hashRange[1]
		}
	}

	if len(t.TagsModel.LSI) >= 5 {
		return errors.New("max local secondary index reached")
	}

	if lsIndex.IndexName != "" && lsIndex.Hash != "" {
		t.TagsModel.LSI = append(t.TagsModel.LSI, *lsIndex)
	}

	return nil
}

// ExtractTypes é um método para extrair os tipos definidos na tag diinamo
// pela chave type
func (t *TagMapper) ExtractTypes(tagsPair []string, field reflect.StructField) error {
	for _, tag := range tagsPair {
		if t.TagsModel.Types == nil {
			t.TagsModel.Types = map[string]reflect.Kind{}
		}

		typeMeta := strings.Split(tag, ":")

		switch typeMeta[0] {
		case _type:
			t.TagsModel.Types[field.Name] = field.Type.Kind()
		}
	}
	return nil
}

// SetPropertyTypes define o valor de PropertyTypes
func (t *TagMapper) SetPropertyTypes(v reflect.Type) {
	t.PropertyTypes = v
}

// GetHash devolve o nome do campo que representa o Hash da tabela
func (t *TagMapper) GetHash() string {
	return t.Hash
}

// GetRange devolve o nome do campo que representa a Sort Key da tabela
func (t *TagMapper) GetRange() string {
	return t.Range
}

// GetType recupera o valor definido pela tag type de uma key específica
func (t *TagMapper) GetType(key string) reflect.Kind {
	return t.Types[key]
}

// GetModel recupera o modelo de tags
func (t *TagMapper) GetModel() *TagsModel {
	return t.TagsModel
}
