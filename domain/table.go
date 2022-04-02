package domain

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	tagManager "github.com/startup-of-zero-reais/dynamo-for-lambda/tag-manager"
)

type (
	// TableClass é o modelo string das classes de tabela
	// STANDARD ou INFREQUENT_ACCESS
	TableClass string

	// Table é o contrato de métodos que uma Table deve
	// implementar para ser possível montar as operações
	// de query
	Table interface {
		GetMetadata() tagManager.Manager
		AttributeDefinitions() []types.AttributeDefinition
		KeySchema() []types.KeySchemaElement
		Billing() types.BillingMode
		ProvisionedThroughput() *types.ProvisionedThroughput
		TableClass() types.TableClass

		GetGSI() []types.GlobalSecondaryIndex
		GetLSI() []types.LocalSecondaryIndex
	}
)
