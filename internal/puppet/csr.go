package puppet

import (
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/micromdm/scep/v2/cryptoutil/x509util"
)

var ErrInvalidPEM = errors.New("failed to decode PEM block")

// GetChallengePassword extracts the challenge password from a PEM-encoded CSR.
func GetChallengePassword(pemBytes []byte) (string, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return "", ErrInvalidPEM
	}

	// https://github.com/golang/go/issues/15995
	// https://github.com/micromdm/scep/pull/45
	//
	// The pem package is not able to parse challenge passwords yet,
	// so we need to obtain that through some parsing of our own.
	csrCP, err := x509util.ParseChallengePassword(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("parsing challenge password: %w", err)
	}

	return csrCP, nil
}
