package store

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type mockDynamoDB struct {
	putItemFn    func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	getItemFn    func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	deleteItemFn func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	scanFn       func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

func (m *mockDynamoDB) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if m.putItemFn != nil {
		return m.putItemFn(ctx, params, optFns...)
	}
	return &dynamodb.PutItemOutput{}, nil
}

func (m *mockDynamoDB) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if m.getItemFn != nil {
		return m.getItemFn(ctx, params, optFns...)
	}
	return &dynamodb.GetItemOutput{}, nil
}

func (m *mockDynamoDB) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	if m.deleteItemFn != nil {
		return m.deleteItemFn(ctx, params, optFns...)
	}
	return &dynamodb.DeleteItemOutput{}, nil
}

func (m *mockDynamoDB) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	if m.scanFn != nil {
		return m.scanFn(ctx, params, optFns...)
	}
	return &dynamodb.ScanOutput{}, nil
}

func TestDynamoDBStore_Create(t *testing.T) {
	mock := &mockDynamoDB{
		putItemFn: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			if aws.ToString(params.TableName) != "test-table" {
				t.Errorf("expected table name test-table, got %s", aws.ToString(params.TableName))
			}
			return &dynamodb.PutItemOutput{}, nil
		},
	}

	s := NewDynamoDBStoreWithAPI(mock, "test-table")
	err := s.Create(context.Background(), 1700000000, "node.example.com", "token123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDynamoDBStore_Read(t *testing.T) {
	expectedToken := OTPToken{
		ExpireAtUnix:   1700000000,
		FQDN:           "node.example.com",
		TokenTableItem: "token123",
	}
	item, _ := attributevalue.MarshalMap(expectedToken)

	t.Run("Success", func(t *testing.T) {
		mock := &mockDynamoDB{
			getItemFn: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{
					Item: item,
				}, nil
			},
		}

		s := NewDynamoDBStoreWithAPI(mock, "test-table")
		token, err := s.Read(context.Background(), "node.example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(token, expectedToken) {
			t.Errorf("got token %+v, want %+v", token, expectedToken)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock := &mockDynamoDB{
			getItemFn: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{Item: nil}, nil
			},
		}

		s := NewDynamoDBStoreWithAPI(mock, "test-table")
		_, err := s.Read(context.Background(), "missing.example.com")
		if !errors.Is(err, ErrTokenNotFound) {
			t.Fatalf("expected ErrTokenNotFound, got %v", err)
		}
	})
}

func TestDynamoDBStore_Delete(t *testing.T) {
	mock := &mockDynamoDB{
		deleteItemFn: func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
			if aws.ToString(params.TableName) != "test-table" {
				t.Errorf("expected table name test-table, got %s", aws.ToString(params.TableName))
			}
			return &dynamodb.DeleteItemOutput{}, nil
		},
	}

	s := NewDynamoDBStoreWithAPI(mock, "test-table")
	err := s.Delete(context.Background(), "node.example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDynamoDBStore_ReadAll(t *testing.T) {
	expectedTokens := []OTPToken{
		{
			ExpireAtUnix:   1700000000,
			FQDN:           "node1.example.com",
			TokenTableItem: "token1",
		},
		{
			ExpireAtUnix:   1700000005,
			FQDN:           "node2.example.com",
			TokenTableItem: "token2",
		},
	}

	var items []map[string]types.AttributeValue
	for _, tok := range expectedTokens {
		item, _ := attributevalue.MarshalMap(tok)
		items = append(items, item)
	}

	mock := &mockDynamoDB{
		scanFn: func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
			return &dynamodb.ScanOutput{
				Items: items,
			}, nil
		},
	}

	s := NewDynamoDBStoreWithAPI(mock, "test-table")
	tokens, err := s.ReadAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("got tokens %+v, want %+v", tokens, expectedTokens)
	}
}
