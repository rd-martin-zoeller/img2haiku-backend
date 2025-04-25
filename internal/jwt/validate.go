package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

func Validate(tokenString string, pubKeyStr string) (bool, error) {
	key, err := loadPublicKey(pubKeyStr)
	if err != nil {
		return false, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if !validateAud(claims) {
			return false, errors.New("invalid audience")
		}
		return true, nil
	}

	return false, errors.New("invalid token")
}

func validateAud(claims jwt.MapClaims) bool {
	aud, ok := claims["aud"].(string)
	if !ok {
		return false
	}
	return aud == "img2haiku-backend"
}
