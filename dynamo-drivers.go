package dynamo_for_lambda

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	_ "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"reflect"
)

type (
	Action      string
	Environment string

	Config struct {
		TableName    string
		HashKeyName  string
		RangeKeyName *string
		Environment  Environment
		Endpoint     string
		*Table
	}

	Dynamo interface {
		Perform(action Action, sql SqlExpression, result interface{}) error
		NewExpressionBuilder() SqlExpression
		Migrate() error
		Seed(items ...*dynamodb.PutItemInput) error
	}

	DynamoClient struct {
		Client   *dynamodb.Client
		Ctx      context.Context
		table    *string
		hashKey  *string
		rangeKey *string

		*Table
	}
)

const (
	GET    = Action("GET")
	PUT    = Action("PUT")
	UPDATE = Action("UPDATE")
	DELETE = Action("DELETE")

	prod = Environment("production")
	stg  = Environment("staging")
	dev  = Environment("development")
)

func (e Environment) isDev() bool {
	return string(e) == "development"
}

func NewDynamoClient(ctx context.Context, conf *Config) *DynamoClient {
	if conf.Environment == "" {
		log.Println("empty environment, setting to development...")
		conf.Environment = dev
	}

	var configs []func(options *config.LoadOptions) error
	if conf.Environment.isDev() {
		if conf.Endpoint == "" {
			conf.Endpoint = "http://localhost:8000"
		}

		configs = devConfigs(conf)
	}

	cfg, err := config.LoadDefaultConfig(ctx, configs...)

	if err != nil {
		log.Fatalf("failed on load config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	dynamoClient := &DynamoClient{
		Client:   client,
		Ctx:      ctx,
		table:    aws.String(conf.TableName),
		hashKey:  aws.String(conf.HashKeyName),
		rangeKey: conf.RangeKeyName,
		Table:    conf.Table,
	}

	if conf.Environment.isDev() {
		err = dynamoClient.Migrate()
		if err != nil {
			log.Fatalf("failed on migrate: %v", err)
		}
	}

	return dynamoClient
}

func (d *DynamoClient) Perform(action Action, sql SqlExpression, target interface{}) error {
	if reflect.TypeOf(target).Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}

	switch action {
	case GET:
		return d.Get(sql, target)
	}
	return nil
}

func (d *DynamoClient) NewExpressionBuilder() SqlExpression {
	return NewSqlBuilder(&Config{
		HashKeyName:  d.Table.tableMetadata["hash"].(string),
		RangeKeyName: aws.String(d.Table.tableMetadata["range"].(string)),
	})
}

func (d *DynamoClient) Migrate() error {
	tables, err := d.Client.ListTables(d.Ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatalf("failed to list tables: %v", err)
	}

	created := false
	for _, table := range tables.TableNames {
		if table == *d.table {
			created = true
			break
		}
	}

	if created {
		return nil
	}

	output, err := d.Client.CreateTable(d.Ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions:  d.Table.AttributeDefinitions(),
		KeySchema:             d.Table.KeySchema(),
		TableName:             d.table,
		BillingMode:           d.Table.Billing,
		ProvisionedThroughput: d.Table.ProvisionedThroughput(),
		TableClass:            types.TableClass(d.Table.TableClass),
	})

	log.Println(output)

	return err
}

func (d *DynamoClient) Seed(items ...map[string]types.AttributeValue) error {
	if len(items) <= 0 {
		log.Println("no items to seed")
		return nil
	}

	log.Printf("seeding table with %d items\n", len(items))

	var transactItems []types.TransactWriteItem

	for _, item := range items {
		transactItems = append(transactItems, types.TransactWriteItem{
			Put: &types.Put{
				Item:      item,
				TableName: d.table,
			},
		})
	}

	_, err := d.Client.TransactWriteItems(d.Ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})

	if err != nil {
		return err
	}

	if len(transactItems) > 0 {
		log.Printf("seed complete: %+v seeded", len(items))
	}

	return nil
}

func (d *DynamoClient) FlushDb() {
	p := dynamodb.NewScanPaginator(d.Client, &dynamodb.ScanInput{TableName: d.table, Limit: aws.Int32(10)})

	for p.HasMorePages() {
		page, err := p.NextPage(d.Ctx)
		if err != nil {
			log.Fatalf("failed on paginate: %v", err)
		}

		for _, item := range page.Items {
			keyCondition := map[string]types.AttributeValue{}

			keyCondition[*d.hashKey] = item[*d.hashKey]
			keyCondition[*d.rangeKey] = item[*d.rangeKey]

			_, err = d.Client.DeleteItem(d.Ctx, &dynamodb.DeleteItemInput{
				TableName: d.table,
				Key:       keyCondition,
			})

			if err != nil {
				log.Fatalf("failed on delete: %v", err)
			}
		}
	}

	log.Println("db flush complete")
}

func devConfigs(conf *Config) []func(options *config.LoadOptions) error {
	var configs []func(options *config.LoadOptions) error

	log.Printf("dev environment, setting dynamo endpoint to %s\n", conf.Endpoint)
	configs = append(configs,
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: conf.Endpoint}, nil
				},
			),
		),
		config.WithCredentialsProvider(
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     "TEST",
					SecretAccessKey: "TEST",
					SessionToken:    "TEST",
					Source:          "Hard-coded credentials; values are irrelevant for local DynamoDB",
				},
			},
		),
		config.WithRegion("us-east-1"),
	)

	return configs
}
