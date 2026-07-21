package store

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const DefaultTableName = "puppet-dynamodb-otp"

// DynamoDBAPI defines the raw DynamoDB SDK operations required by DynamoDBStore.
type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

// DynamoDBStore implements Store interface backed by AWS DynamoDB.
type DynamoDBStore struct {
	client    DynamoDBAPI
	tableName string
}

// GetDefaultTableName returns table name from DYNAMODB_TABLE_NAME or default.
func GetDefaultTableName() string {
	if tableName := os.Getenv("DYNAMODB_TABLE_NAME"); tableName != "" {
		return tableName
	}
	return DefaultTableName
}

// NewDynamoDBStore creates a new DynamoDBStore with AWS default configuration.
func NewDynamoDBStore(ctx context.Context) (*DynamoDBStore, error) {
	return NewDynamoDBStoreWithTableName(ctx, GetDefaultTableName())
}

// NewDynamoDBStoreWithTableName creates a new DynamoDBStore with a custom table name.
func NewDynamoDBStoreWithTableName(ctx context.Context, tableName string) (*DynamoDBStore, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	return NewDynamoDBStoreWithAPI(dynamodb.NewFromConfig(cfg), tableName), nil
}

// NewDynamoDBStoreWithAPI constructs a DynamoDBStore with a custom DynamoDBAPI client.
func NewDynamoDBStoreWithAPI(api DynamoDBAPI, tableName string) *DynamoDBStore {
	if tableName == "" {
		tableName = DefaultTableName
	}
	return &DynamoDBStore{
		client:    api,
		tableName: tableName,
	}
}

// Create stores a new OTP token in DynamoDB.
func (s *DynamoDBStore) Create(ctx context.Context, expireAtUnix int64, fqdn string, otpToken string) error {
	token := OTPToken{
		ExpireAtUnix:   expireAtUnix,
		FQDN:           fqdn,
		TokenTableItem: otpToken,
	}

	item, err := attributevalue.MarshalMap(token)
	if err != nil {
		return fmt.Errorf("marshaling token attribute values: %w", err)
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("creating OTP token for %s: %w", fqdn, err)
	}

	return nil
}

// Delete removes an OTP token from DynamoDB by FQDN.
func (s *DynamoDBStore) Delete(ctx context.Context, fqdn string) error {
	_, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"fqdn": &types.AttributeValueMemberS{Value: fqdn},
		},
	})
	if err != nil {
		return fmt.Errorf("deleting OTP token for %s: %w", fqdn, err)
	}

	return nil
}

// Read fetches an OTP token from DynamoDB by FQDN.
func (s *DynamoDBStore) Read(ctx context.Context, fqdn string) (OTPToken, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"fqdn": &types.AttributeValueMemberS{Value: fqdn},
		},
	})
	if err != nil {
		return OTPToken{}, fmt.Errorf("reading OTP token for %s: %w", fqdn, err)
	}

	if len(out.Item) == 0 {
		return OTPToken{}, fmt.Errorf("%w for %s", ErrTokenNotFound, fqdn)
	}

	var token OTPToken
	err = attributevalue.UnmarshalMap(out.Item, &token)
	if err != nil {
		return OTPToken{}, fmt.Errorf("unmarshaling OTP token for %s: %w", fqdn, err)
	}

	return token, nil
}

// ReadAll scans and retrieves all OTP tokens from DynamoDB.
func (s *DynamoDBStore) ReadAll(ctx context.Context) ([]OTPToken, error) {
	out, err := s.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(s.tableName),
	})
	if err != nil {
		return nil, fmt.Errorf("scanning OTP tokens: %w", err)
	}

	var tokens []OTPToken
	err = attributevalue.UnmarshalListOfMaps(out.Items, &tokens)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling OTP tokens: %w", err)
	}

	return tokens, nil
}

// Verify interface compliance at compile time.
var _ Store = (*DynamoDBStore)(nil)
