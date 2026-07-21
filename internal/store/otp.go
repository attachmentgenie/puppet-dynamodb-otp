package store

import (
	"context"
	"errors"
)

var ErrTokenNotFound = errors.New("token not found")

// OTPToken represents an OTP record.
type OTPToken struct {
	ExpireAtUnix   int64  `dynamodbav:"expire_at_unix"`
	FQDN           string `dynamodbav:"fqdn"`
	TokenTableItem string `dynamodbav:"token_table_item"`
}

// Store defines the domain interface for OTP token operations.
type Store interface {
	Create(ctx context.Context, expireAtUnix int64, fqdn string, otpToken string) error
	Delete(ctx context.Context, fqdn string) error
	Read(ctx context.Context, fqdn string) (OTPToken, error)
	ReadAll(ctx context.Context) ([]OTPToken, error)
}
