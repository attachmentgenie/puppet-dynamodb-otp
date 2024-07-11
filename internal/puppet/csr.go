package puppet

import (
	"encoding/pem"
	"log"

	"github.com/micromdm/scep/v2/cryptoutil/x509util"
)

func GetChallengePassword(PEM []byte) (string, error) {
	block, _ := pem.Decode(PEM)
	if block == nil {
		log.Fatal("failed to decode PEM block")
	}
	// https://github.com/golang/go/issues/15995
	// https://github.com/micromdm/scep/pull/45
	//
	// The pem package is not able to parse challenge passwords yet,
	// so we need to obtain that through some parsing of our own.
	// Luckily someone did the hard work for us already
	csrCP, err := x509util.ParseChallengePassword(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	return csrCP, err
}
