package main

import (
	"log"

	"github.com/lapitskyss/chat-service/internal/starter"
)

func main() {
	srv, cleanup, err := starter.Initialize()
	if err != nil {
		log.Fatalf("initialize app: %v", err)
	}

	if err := srv.Run(cleanup); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
