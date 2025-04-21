package main

import (
	"log"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"

	// Blank-import the function package so the init() runs
	_ "github.com/rd-martin-zoeller/img2haiku-backend"
)

const (
	sub = "img2haiku-backend-demo"
	aud = "img2haiku-backend"
	ttl = 15 * time.Minute
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

	prepJWT()

	log.Printf("Starting function framework on %s:%s\n", hostname, port)

	if err := funcframework.StartHostPort(hostname, port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}
