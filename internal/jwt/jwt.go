package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	privateKeyEnvVar = "JWT_PRIVATE_KEY"
	publicKeyEnvVar  = "JWT_PUBLIC_KEY"
	sub              = "img2haiku-backend-demo"
	aud              = "img2haiku-backend"
	ttl              = 15 * time.Minute
)

func JWTForTesting() (string, error) {
	key, err := loadPrivateKey()
	if err != nil {
		return "", err
	}

	now := time.Now()
	jwt := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": sub,
		"aud": aud,
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
	})

	jwtString, err := jwt.SignedString(key)
	if err != nil {
		return "", err
	}

	// Validate the JWT
	valid, err := Validate(jwtString)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", errors.New("invalid JWT")
	}

	return jwtString, nil
}
