package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

const privateKeyEnvVar = "JWT_PRIVATE_KEY"

func jwtForTesting() (string, error) {
	_, err := loadPrivateKey()
	if err != nil {
		return "", err
	}

	return "placeholder-jwt-for-testing", nil
}

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
