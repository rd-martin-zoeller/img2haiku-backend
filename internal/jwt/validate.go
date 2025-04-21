package jwt

import (
	"errors"
	"fmt"

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
		fmt.Println("Claims:", claims)

		if !validateAud(claims) {
			return false, errors.New("invalid audience")
		}
		if !validateSub(claims) {
			return false, errors.New("invalid subject")
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

func validateSub(claims jwt.MapClaims) bool {
	sub, ok := claims["sub"].(string)
	if !ok {
		return false
	}
	return sub == "img2haiku-backend-demo"
}
