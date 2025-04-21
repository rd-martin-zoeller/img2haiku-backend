package main

import (
	"log"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/rd-martin-zoeller/img2haiku-backend/internal/jwt"

	// Blank-import the function package so the init() runs
	_ "github.com/rd-martin-zoeller/img2haiku-backend"
)

const (
	sub = "img2haiku-backend-demo"
	aud = "img2haiku-backend"
	ttl = 15 * time.Minute
)

const (
	privateKeyEnvVar = "JWT_PRIVATE_KEY"
	publicKeyEnvVar  = "JWT_PUBLIC_KEY"
)

func main() {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	}

	jwt, err := jwt.JWTForTesting(jwt.JWTConfig{
		KeyPair: jwt.KeyPair{
			Private: os.Getenv(privateKeyEnvVar),
			Public:  os.Getenv(publicKeyEnvVar),
		},
		Sub: sub,
		Aud: aud,
		Exp: ttl,
	})
	if err != nil {
		log.Fatalf("Could not generate JWT for testing: %v\n", err)
	}

	if jwt == "" {
		log.Fatalf("JWT for testing is empty\n")
	}
	log.Printf("JWT for testing: %s\n", jwt)

	log.Printf("Starting function framework on %s:%s\n", hostname, port)

	if err := funcframework.StartHostPort(hostname, port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}
