package domain

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/startup-of-zero-reais/dynamo-for-lambda/logger"
)

type (
	Config struct {
		TableName   string
		Environment Environment
		Client      *dynamodb.Client
		Table
		logger.Log
	}
)
