package puppet

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"testing"
)

type rawAttribute struct {
	Type  asn1.ObjectIdentifier
	Value asn1.RawValue `asn1:"set"`
}

type certificationRequestInfo struct {
	Version int
	Subject pkix.RDNSequence
	PKI     asn1.RawValue
	Attributes []rawAttribute `asn1:"tag:0,optional"`
}

type certificationRequest struct {
	Info      certificationRequestInfo
	Algorithm pkix.AlgorithmIdentifier
	Signature asn1.BitString
}

func generateTestCSR(t *testing.T, challengePassword string) []byte {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}

	var pubInfo struct {
		Algorithm pkix.AlgorithmIdentifier
		PublicKey asn1.BitString
	}
	_, err = asn1.Unmarshal(pubBytes, &pubInfo)
	if err != nil {
		t.Fatalf("failed to unmarshal public key info: %v", err)
	}

	oidChallengePassword := asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 7}

	var attributes []rawAttribute
	if challengePassword != "" {
		passBytes, err := asn1.Marshal(asn1.RawValue{
			Tag:   asn1.TagPrintableString,
			Class: asn1.ClassUniversal,
			Bytes: []byte(challengePassword),
		})
		if err != nil {
			t.Fatalf("failed to marshal passBytes: %v", err)
		}

		attributes = append(attributes, rawAttribute{
			Type: oidChallengePassword,
			Value: asn1.RawValue{
				Tag:        asn1.TagSet,
				Class:      asn1.ClassUniversal,
				IsCompound: true,
				Bytes:      passBytes,
			},
		})
	}

	info := certificationRequestInfo{
		Version: 0,
		Subject: pkix.RDNSequence{
			{
				{
					Type:  asn1.ObjectIdentifier{2, 5, 4, 3}, // Common Name
					Value: "test-node.example.com",
				},
			},
		},
		PKI: asn1.RawValue{
			FullBytes: pubBytes,
		},
		Attributes: attributes,
	}

	infoBytes, err := asn1.Marshal(info)
	if err != nil {
		t.Fatalf("failed to marshal info: %v", err)
	}

	// OID for SHA256 with RSA
	sigAlg := pkix.AlgorithmIdentifier{
		Algorithm: asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 11},
	}

	csr := certificationRequest{
		Info:      info,
		Algorithm: sigAlg,
		Signature: asn1.BitString{
			Bytes:     infoBytes, // dummy signature for test parsing
			BitLength: len(infoBytes) * 8,
		},
	}

	csrBytes, err := asn1.Marshal(csr)
	if err != nil {
		t.Fatalf("failed to marshal csr: %v", err)
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})

	return csrPEM
}

func TestGetChallengePassword(t *testing.T) {
	tests := []struct {
		name          string
		pemInput      []byte
		wantPassword  string
		wantErr       bool
		expectedError string
	}{
		{
			name:          "Invalid PEM Block",
			pemInput:      []byte("invalid pem data"),
			wantErr:       true,
			expectedError: ErrInvalidPEM.Error(),
		},
		{
			name:         "Valid CSR with Challenge Password",
			pemInput:     generateTestCSR(t, "secret-otp-token-12345"),
			wantPassword: "secret-otp-token-12345",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetChallengePassword(tt.pemInput)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetChallengePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.expectedError != "" {
				if err.Error() != tt.expectedError {
					t.Errorf("GetChallengePassword() error = %v, want %v", err.Error(), tt.expectedError)
				}
			}
			if !tt.wantErr && got != tt.wantPassword {
				t.Errorf("GetChallengePassword() = %v, want %v", got, tt.wantPassword)
			}
		})
	}
}
