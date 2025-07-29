package ams

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func TestAuthorizationManagerforIAS(t *testing.T) {
	t.Run("with broken cert", func(t *testing.T) {
		_, err := NewAuthorizationManagerForIAS("https://example.com", "brokencert", "test", "test")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("with broken url", func(t *testing.T) {
		// create simple valid cert
		cert, key := generateTestCert(t)
		_, err := NewAuthorizationManagerForIAS("noprot://example.com ", "dummy-id", string(cert), string(key))
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("with valid cert", func(t *testing.T) {
		// create simple valid cert
		cert, key := generateTestCert(t)
		_, err := NewAuthorizationManagerForIAS("https://example.com", "dummy-id", string(cert), string(key))
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestAuthorizationManagerforLocal(t *testing.T) {
	t.Run("with broken cert", func(t *testing.T) {
		a := NewAuthorizationManagerForFs("/tmp")

		if a == nil {
			t.Errorf("Expected non-nil, got nil")
		}
	})
}

func generateTestCert(t *testing.T) ([]byte, []byte) {
	t.Helper()
	var certPEM, keyPEM []byte

	// Generate private key
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	// Create a template for the certificate
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour * 24), // Valid for 1 day

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create the certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	// Encode certificate to PEM
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// Encode private key to PEM
	keyBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		t.Fatalf("failed to marshal private key: %v", err)
	}
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})

	return certPEM, keyPEM
}
