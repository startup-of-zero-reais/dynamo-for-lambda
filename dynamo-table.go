package dynamo_for_lambda

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"reflect"
	"strings"
)

type (
	TableClass string

	Table struct {
		TableName  string
		Billing    types.BillingMode
		TableClass TableClass
		Keys       interface{}

		ReadThroughput  int32
		WriteThroughput int32

		tableMetadata TableMetadata
	}

	TableMetadata map[string]interface{}

	Mocktable struct {
		PK           string `diinamo:"type:string;hash"`
		SK           string `diinamo:"type:string;range"`
		Owner        string `diinamo:"type:string;gsi:CourseOwnerIndex;keyPairs:PK=Owner"`
		Title        string `diinamo:"type:string;gsi:CourseTitleIndex;keyPairs:Title=SK"`
		ParentCourse string `diinamo:"type:string;gsi:CourseLessonsIndex;keyPairs:ParentCourse=SK"`
		ParentModule string `diinamo:"type:string;gsi:ModuleLessonsIndex;keyPairs:ParentModule=SK"`
	}
)

const (
	STANDARD          = TableClass("STANDARD")
	INFREQUEST_ACCESS = TableClass("INFREQUEST_ACCESS")
)

func NewTable(tbName string, tableStruct interface{}) *Table {
	if reflect.TypeOf(tableStruct).Kind() != reflect.Struct {
		log.Fatalf("table struct should be a struct")
	}

	t := &Table{
		TableName:       tbName,
		Billing:         types.BillingModeProvisioned,
		TableClass:      STANDARD,
		Keys:            tableStruct,
		ReadThroughput:  int32(1),
		WriteThroughput: int32(1),
	}

	t.extractTagMap()

	return t
}

func (t *Table) AttributeDefinitions() []types.AttributeDefinition {
	return []types.AttributeDefinition{
		t.getAttrDefinition(t.tableMetadata["hash"].(string)),
		t.getAttrDefinition(t.tableMetadata["range"].(string)),
	}
}

func (t *Table) KeySchema() []types.KeySchemaElement {
	return t.getKeySchema()
}

func (t *Table) getAttrDefinition(key string) types.AttributeDefinition {
	mapHashType := t.tableMetadata["type"].(map[string]string)[key]

	hashType := types.ScalarAttributeTypeS
	switch mapHashType {
	case "number":
		hashType = types.ScalarAttributeTypeN
	case "binary":
		hashType = types.ScalarAttributeTypeB
	}

	return types.AttributeDefinition{
		AttributeName: aws.String(t.tableMetadata["hash"].(string)),
		AttributeType: hashType,
	}
}

func (t *Table) getKeySchema() []types.KeySchemaElement {
	var keySchema []types.KeySchemaElement

	if t.tableMetadata["hash"] != nil {
		keySchema = append(keySchema, types.KeySchemaElement{
			AttributeName: aws.String(t.tableMetadata["hash"].(string)),
			KeyType:       types.KeyTypeHash,
		})
	}

	if t.tableMetadata["range"] != nil {
		keySchema = append(keySchema, types.KeySchemaElement{
			AttributeName: aws.String(t.tableMetadata["range"].(string)),
			KeyType:       types.KeyTypeRange,
		})
	}

	return keySchema
}

func (t *Table) ProvisionedThroughput() *types.ProvisionedThroughput {
	return &types.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(int64(t.ReadThroughput)),
		WriteCapacityUnits: aws.Int64(int64(t.WriteThroughput)),
	}
}

func (t *Table) extractTagMap() {
	keyTypes := reflect.TypeOf(t.Keys)

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

	t.tableMetadata = tagMap
}

func GenTable() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{},
		KeySchema:            []types.KeySchemaElement{},
		TableName:            aws.String(""),
		BillingMode:          types.BillingModeProvisioned,
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String(""),
				KeySchema: []types.KeySchemaElement{},
				Projection: &types.Projection{
					NonKeyAttributes: nil,
					ProjectionType:   "ALL",
				},
				ProvisionedThroughput: nil,
			},
		},
		LocalSecondaryIndexes: []types.LocalSecondaryIndex{},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  nil,
			WriteCapacityUnits: nil,
		},
		TableClass: "STANDARD",
	}
}
