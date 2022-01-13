package tagManager

import (
	"reflect"
	"strings"
)

type (
	// TagManager é a estrutura que gerencia as tags extraídas.
	// Tags específicas da chave: diinamo
	TagManager struct {
		StructToMap interface{}
		Tags        map[string]interface{}

		Log

		*TagMapper
	}
)

// NewTagManager inicializa uma estrutura TagManager
func NewTagManager() *TagManager {
	logger := NewLogger()

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

func (t *TagManager) extractTagMap() {
	keyTypes := reflect.TypeOf(t.PropertyTypes)

	tagMap := map[string]interface{}{}

	for i := 0; i < keyTypes.NumField(); i++ {
		field := keyTypes.Field(i)

		if v, ok := field.Tag.Lookup("diinamo"); ok {
			tag := strings.Split(v, ";")
			for _, props := range tag {
				kv := strings.Split(props, ":")

				if len(kv) >= 1 {
					var value string
					if len(kv) > 1 {
						value = kv[1]
					}

					key := kv[0]

					switch key {
					case "hash", "range":
						tagMap[key] = field.Name
					case "type":
						if s := tagMap[key]; s == nil {
							tagMap[key] = map[string]string{}
						}

						tagMap[key].(map[string]string)[field.Name] = value
					case "gsi":
						if s := tagMap[key]; s == nil {
							tagMap[key] = map[string]map[string]string{}
						}

						tagMap[key].(map[string]map[string]string)[value] = map[string]string{}
					case "keyPairs":
						ikv := strings.Split(value, "=")
						ihash := ikv[0]
						irang := ikv[1]

						for index, _ := range tagMap["gsi"].(map[string]map[string]string) {
							tagMap["gsi"].(map[string]map[string]string)[index]["hash"] = ihash
							tagMap["gsi"].(map[string]map[string]string)[index]["range"] = irang
						}
					}
				}
			}
		}
	}

	//t.tableMetadata = tagMap
}
