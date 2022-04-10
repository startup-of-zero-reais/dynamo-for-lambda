package domain

import "github.com/startup-of-zero-reais/dynamo-for-lambda/logger"

type (
	Config struct {
		TableName   string
		Environment Environment
		Endpoint    string
		Region      string
		Table
		logger.Log
	}
)
