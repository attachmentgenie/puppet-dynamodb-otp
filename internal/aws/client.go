package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Client struct {
	dynamodb *dynamodb.Client
}

func New() (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("ERROR %s", err)
	}
	client := &Client{
		dynamodb: dynamodb.NewFromConfig(cfg),
	}
	return client, nil
}
