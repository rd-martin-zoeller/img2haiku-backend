package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func loadPrivateKey(keyStr string) (*rsa.PrivateKey, error) {
	if keyStr == "" {
		return nil, errors.New("keyStr var not set: " + keyStr)
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

func loadPublicKey(keyStr string) (*rsa.PublicKey, error) {
	if keyStr == "" {
		return nil, errors.New("keyStr var not set: " + keyStr)
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
