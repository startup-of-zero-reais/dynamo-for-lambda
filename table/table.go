package table

import (
	"log"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/drivers"
	tagManager "github.com/startup-of-zero-reais/dynamo-for-lambda/tag-manager"
)

type (
	// Table é a estrutura de tabela necessária para criar a tabela no DynamoDB
	Table struct {
		TableName      string
		BillingMode    types.BillingMode
		TableClassMode domain.TableClass

		ReadThroughput  int32
		WriteThroughput int32

		Metadata tagManager.Manager
	}
)

// NewTable é um construtor e inicializador para uma estrutura Table
func NewTable(tbName string, tableStruct interface{}) *Table {
	if reflect.TypeOf(tableStruct).Kind() != reflect.Struct {
		log.Fatalf("table struct should be a struct")
	}

	t := &Table{
		TableName:       tbName,
		BillingMode:     types.BillingModeProvisioned,
		TableClassMode:  drivers.STANDARD,
		ReadThroughput:  int32(1),
		WriteThroughput: int32(1),
	}

	t.Metadata = tagManager.NewTagManager().SetEntity(tableStruct)
	err := t.Metadata.MapTags()
	if err != nil {
		log.Fatalf("error on map tags")
	}

	return t
}

// GetMetadata é o método para retornar o Manager dentro da tabela
func (t *Table) GetMetadata() tagManager.Manager {
	return t.Metadata
}

// AttributeDefinitions é o método que retorna a definição de atributos para a PK.
// A PK é composta de Hash e Range keys
func (t *Table) AttributeDefinitions() []types.AttributeDefinition {
	attrDefinitions := []types.AttributeDefinition{
		t.getAttrDefinition(t.Metadata.GetHash()),
		t.getAttrDefinition(t.Metadata.GetRange()),
	}

	log.Println("performing attributes definition")
	// Adiciona os atributos de Global Secondary Index às definições
	// de atributos da Tabela
	for _, gsi := range t.GetGSI() {
		for _, key := range gsi.KeySchema {
			issetField := false

			for _, attr := range attrDefinitions {
				log.Printf("%s == %s\n", *attr.AttributeName, *key.AttributeName)
				if *attr.AttributeName == *key.AttributeName && !issetField {
					issetField = true
				}
			}

			// Se o atributo não estiver na lista de atributos, adiciona
			if !issetField {
				attrDefinitions = append(attrDefinitions, t.getAttrDefinition(*key.AttributeName))
			}
		}
	}

	log.Printf("attr definitions: %v", attrDefinitions)

	return attrDefinitions
}

// KeySchema é o método que retorna o schema de chaves que compõe a PK
func (t *Table) KeySchema() []types.KeySchemaElement {
	return t.getKeySchema()
}

// getAttrDefinition é um método que retorna o types.AttributeDefinition de uma
// chave específica contida na entity de Metadata
func (t *Table) getAttrDefinition(key string) types.AttributeDefinition {
	mapHashType := t.Metadata.GetType(key)

	if mapHashType == reflect.Invalid {
		log.Fatalf("invalid key attr definition")
	}

	hashType := types.ScalarAttributeTypeS
	switch mapHashType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		hashType = types.ScalarAttributeTypeN
	case reflect.Map, reflect.Slice, reflect.Array:
		hashType = types.ScalarAttributeTypeB
	}

	return types.AttributeDefinition{
		AttributeName: aws.String(key),
		AttributeType: hashType,
	}
}

// getKeySchema é o método que monta o schema da chave PK
// utilizando o Hash e a Range keys
func (t *Table) getKeySchema() []types.KeySchemaElement {
	var keySchema []types.KeySchemaElement

	if t.Metadata.GetHash() != "" {
		keySchema = append(keySchema, types.KeySchemaElement{
			AttributeName: aws.String(t.Metadata.GetHash()),
			KeyType:       types.KeyTypeHash,
		})
	}

	if t.Metadata.GetRange() != "" {
		keySchema = append(keySchema, types.KeySchemaElement{
			AttributeName: aws.String(t.Metadata.GetRange()),
			KeyType:       types.KeyTypeRange,
		})
	}

	return keySchema
}

// Billing é o método que retorna o BillingMode a ser configurado na tabela
func (t *Table) Billing() types.BillingMode {
	return t.BillingMode
}

// ProvisionedThroughput é o método que retorna o ponteiro de configuração
// do Throughput da tabela.
//
// Os padrões de provisionamento são 1 Read Capacity e 1 Write Capacity
func (t *Table) ProvisionedThroughput() *types.ProvisionedThroughput {
	return &types.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(int64(t.ReadThroughput)),
		WriteCapacityUnits: aws.Int64(int64(t.WriteThroughput)),
	}
}

// TableClass é o método que retorna o TableClass da tabela
func (t *Table) TableClass() types.TableClass {
	return types.TableClass(t.TableClassMode)
}

// GetGSI é o método que monta e retorna os GlobalSecondaryIndex
func (t *Table) GetGSI() []types.GlobalSecondaryIndex {
	var GSIs []types.GlobalSecondaryIndex

	model := t.GetMetadata().GetMapper().GetModel()

	for _, gsi := range model.GSI {
		if gsi.IndexName == "" {
			continue
		}

		parseGSI := types.GlobalSecondaryIndex{
			IndexName: aws.String(gsi.IndexName),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String(gsi.Hash),
					KeyType:       types.KeyTypeHash,
				},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(int64(gsi.ProvisionedThroughput.ReadCapacity)),
				WriteCapacityUnits: aws.Int64(int64(gsi.ProvisionedThroughput.WriteCapacity)),
			},
			Projection: &types.Projection{
				ProjectionType: types.ProjectionType("ALL"),
			},
		}

		if gsi.Range != "" {
			parseGSI.KeySchema = append(parseGSI.KeySchema, types.KeySchemaElement{
				AttributeName: aws.String(gsi.Range),
				KeyType:       types.KeyTypeRange,
			})
		}

		GSIs = append(GSIs, parseGSI)
	}

	return GSIs
}

// GetLSI é o método que monta e retorna os LocalSecondaryIndex
func (t *Table) GetLSI() []types.LocalSecondaryIndex {
	var LSIs []types.LocalSecondaryIndex

	model := t.GetMetadata().GetMapper().GetModel()

	for _, lsi := range model.LSI {
		parseLSI := types.LocalSecondaryIndex{
			IndexName: aws.String(lsi.IndexName),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String(lsi.Hash),
					KeyType:       types.KeyTypeHash,
				},
			},
			Projection: &types.Projection{
				ProjectionType: types.ProjectionType("ALL"),
			},
		}

		if lsi.Range != "" {
			parseLSI.KeySchema = append(parseLSI.KeySchema, types.KeySchemaElement{
				AttributeName: aws.String(lsi.Range),
				KeyType:       types.KeyTypeRange,
			})
		}

		LSIs = append(LSIs, parseLSI)
	}

	return LSIs
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
		LocalSecondaryIndexes: []types.LocalSecondaryIndex{
			{
				IndexName:  nil,
				KeySchema:  nil,
				Projection: nil,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  nil,
			WriteCapacityUnits: nil,
		},
		TableClass: "STANDARD",
	}
}
