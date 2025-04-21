package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func loadPrivateKey() (*rsa.PrivateKey, error) {
	keyStr := os.Getenv(privateKeyEnvVar)
	if keyStr == "" {
		return nil, errors.New("env var not set: " + privateKeyEnvVar)
	}

	keyData := []byte(keyStr)
	block, _ := pem.Decode(keyData)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}
	return rsaKey, nil
}

func loadPublicKey() (*rsa.PublicKey, error) {
	keyStr := os.Getenv(publicKeyEnvVar)
	if keyStr == "" {
		return nil, errors.New("env var not set: " + publicKeyEnvVar)
	}

	block, _ := pem.Decode([]byte(keyStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		if rsaPub, ok := pub.(*rsa.PublicKey); ok {
			return rsaPub, nil
		}
		return nil, errors.New("not an RSA public key")
	}

	return nil, errors.New("unsupported public key format")
}
