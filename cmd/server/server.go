package main

import (
	"flag"
	"log"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"

	// Blank-import the function package so the init() runs
	_ "github.com/rd-martin-zoeller/img2haiku-backend"
)

const (
	sub      = "img2haiku-backend-demo"
	aud      = "img2haiku-backend"
	ttl      = 15 * time.Minute
	hostname = "127.0.0.1"
)

var (
	port = flag.String("port", "8080", "Port to run the server on")
)

func main() {
	flag.Parse()

	prepJWT()

	log.Printf("Starting function framework on %s:%s\n", hostname, *port)

	if err := funcframework.StartHostPort(hostname, *port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}
