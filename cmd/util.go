package main

import (
	"log"
	"os"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/jwt"
)

func prepJWT() {
	keyPair, err := jwt.GenKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate credentials: %v\n", err)
	}

	os.Setenv("JWT_SECRET", keyPair.Public)
	jwt, err := jwt.JWTForTesting(jwt.JWTConfig{
		KeyPair: keyPair,
		Sub:     sub,
		Aud:     aud,
		Exp:     ttl,
	})
	if err != nil {
		log.Fatalf("Could not generate JWT for testing: %v\n", err)
	}

	if jwt == "" {
		log.Fatalf("JWT for testing is empty\n")
	}

	log.Printf("+++++ JWT for testing+++++\n%s\n", jwt)
}
