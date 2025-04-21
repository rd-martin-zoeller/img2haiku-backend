package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	sub = "img2haiku-backend-demo"
	aud = "img2haiku-backend"
	ttl = 15 * time.Minute
)

type JWTConfig struct {
	KeyPair KeyPair
	Sub     string
	Aud     string
	Exp     time.Duration
}

type KeyPair struct {
	Private string
	Public  string
}

func JWTForTesting(config JWTConfig) (string, error) {
	key, err := loadPrivateKey(config.KeyPair.Private)
	if err != nil {
		return "", err
	}

	now := time.Now()
	jwt := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": config.Sub,
		"aud": config.Aud,
		"iat": now.Unix(),
		"exp": now.Add(config.Exp).Unix(),
	})

	jwtString, err := jwt.SignedString(key)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}
