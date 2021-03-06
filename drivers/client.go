package drivers

import (
	"context"
	"errors"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/domain"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/expressions"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/logger"
)

type (
	DynamoClient struct {
		Client    *dynamodb.Client
		Ctx       context.Context
		TableName *string
		HashKey   *string
		RangeKey  *string

		domain.Table
		logger.Log
	}
)

func NewDynamoClient(ctx context.Context, conf *domain.Config) *DynamoClient {
	if conf.Log == nil {
		conf.Log = logger.NewLogger()
	}

	dynamoClient := &DynamoClient{
		Client:    conf.Client,
		Ctx:       ctx,
		TableName: aws.String(conf.TableName),
		HashKey:   aws.String(conf.GetMetadata().GetHash()),
		RangeKey:  aws.String(conf.GetMetadata().GetHash()),
		Table:     conf.Table,
		Log:       conf.Log,
	}

	conf.Log.Info("dynamo client connected\n")

	return dynamoClient
}

func (d *DynamoClient) Perform(action domain.Action, sql domain.SqlExpression, target interface{}) error {
	d.Log.Info("performing %s action\n", action)
	if reflect.TypeOf(target).Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}

	switch action {
	case GET:
		return d.Get(sql, target)
	case PUT:
		return d.Put(sql, target)
	case QUERY:
		return d.Query(sql, target)
	case UPDATE:
		return d.Update(sql, target)
	case DELETE:
		return d.Delete(sql)
	}
	return nil
}

func (d *DynamoClient) NewExpressionBuilder() domain.SqlExpression {
	return expressions.NewSqlBuilder(&domain.Config{
		TableName: *d.TableName,
		Table:     d.Table,
		Log:       d.Log,
	})
}

func (d *DynamoClient) Migrate() error {
	tables, err := d.Client.ListTables(d.Ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		d.Critical("failed to list tables: %v", err)
	}

	d.Debug("found %d tables\n", len(tables.TableNames))

	created := false
	for _, table := range tables.TableNames {
		d.Info("checking table: %s\n", table)
		if table == *d.TableName {
			created = true
			break
		}
	}

	if created {
		d.Info("table `%s` already exists\n", *d.TableName)
		return nil
	}

	err = d.CreateTable()
	if err != nil {
		return err
	}

	return err
}

func (d *DynamoClient) Seed(items ...map[string]types.AttributeValue) error {
	if len(items) <= 0 {
		d.Info("no items to seed")
		return nil
	}

	d.Debug("seeding table with %d items\n", len(items))

	var transactItems []types.TransactWriteItem

	for _, item := range items {
		d.Debug("seeding item: %+v\n", item)
		transactItems = append(transactItems, types.TransactWriteItem{
			Put: &types.Put{
				Item:      item,
				TableName: d.TableName,
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
		d.Log.Debug("seed complete: %+v seeded\n", len(items))
	}

	return nil
}

func (d *DynamoClient) FlushDb() {
	d.Error("performing flush db action. Remove this instruction to not lose all your base")
	p := dynamodb.NewScanPaginator(d.Client, &dynamodb.ScanInput{TableName: d.TableName, Limit: aws.Int32(10)})

	for p.HasMorePages() {
		page, err := p.NextPage(d.Ctx)
		if err != nil {
			d.Critical("failed on paginate: %v\n", err)
		}

		for _, item := range page.Items {
			keyCondition := map[string]types.AttributeValue{}

			keyCondition[d.GetMetadata().GetHash()] = item[d.GetMetadata().GetHash()]
			keyCondition[d.GetMetadata().GetRange()] = item[d.GetMetadata().GetRange()]

			_, err = d.Client.DeleteItem(d.Ctx, &dynamodb.DeleteItemInput{
				TableName: d.TableName,
				Key:       keyCondition,
			})

			if err != nil {
				d.Critical("failed on delete: %v\n", err)
			}
		}
	}

	out, err := d.Client.DeleteTable(d.Ctx, &dynamodb.DeleteTableInput{
		TableName: d.TableName,
	})

	if err != nil {
		d.Critical("failed on delete table: %v\n", err)
	}

	d.Info("table '%s' deleted\n", *out.TableDescription.TableName)
	d.Info("db flush complete")
}

func (d *DynamoClient) CreateTable() error {
	table := &dynamodb.CreateTableInput{
		AttributeDefinitions:   d.AttributeDefinitions(),
		KeySchema:              d.KeySchema(),
		TableName:              d.TableName,
		BillingMode:            d.Billing(),
		GlobalSecondaryIndexes: d.GetGSI(),
		LocalSecondaryIndexes:  d.GetLSI(),
		ProvisionedThroughput:  d.ProvisionedThroughput(),
		TableClass:             types.TableClass(d.Table.TableClass()),
	}

	_, err := d.Client.CreateTable(d.Ctx, table)
	if err != nil {
		return err
	}

	d.Debug("table `%+v` created\n", *d.TableName)

	return nil
}
