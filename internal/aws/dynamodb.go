package aws

import (
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var TableName = "puppet-dynamodb-otp"

type Otp_token struct {
	Expire_at_unix int64
	Fqdn           string
	Otp_token      string
}

func Create(expire_at_unix int64, fqdn string, otp_token string) {
	svc := GetDynamodbClient()
	_, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item: map[string]types.AttributeValue{
			"expire_at_unix": &types.AttributeValueMemberN{Value: strconv.FormatInt(expire_at_unix, 10)},
			"fqdn":           &types.AttributeValueMemberS{Value: fqdn},
			"otp_token":      &types.AttributeValueMemberS{Value: otp_token},
		},
	})
	if err != nil {
		panic(err)
	}
}

func Delete(fqdn string) {
	svc := GetDynamodbClient()
	_, err := svc.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"fqdn": &types.AttributeValueMemberS{Value: fqdn},
		},
	})
	if err != nil {
		panic(err)
	}
}

func Read(fqdn string) (Otp_token, error) {
	svc := GetDynamodbClient()
	out, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"fqdn": &types.AttributeValueMemberS{Value: fqdn},
		},
	})

	var token Otp_token
	if err != nil {
		panic(err)
	} else {
		err = attributevalue.UnmarshalMap(out.Item, &token)
		if err != nil {
			log.Fatalf("unmarshal failed, %v", err)
		}
	}
	return token, err
}

func ReadAll() []Otp_token {
	svc := GetDynamodbClient()
	out, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(TableName),
	})
	if err != nil {
		panic(err)
	}

	var tokens []Otp_token
	err = attributevalue.UnmarshalListOfMaps(out.Items, &tokens)
	if err != nil {
		log.Fatalf("unable to unmarshal tokens: %v", err)
	}
	return tokens
}
