package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/rd-martin-zoeller/img2haiku-backend/internal/jwt"

	// Blank-import the function package so the init() runs
	_ "github.com/rd-martin-zoeller/img2haiku-backend"
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

	jwt, err := jwt.JWTForTesting()
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
