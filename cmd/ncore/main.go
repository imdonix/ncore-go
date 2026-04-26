package main

import (
	"log"

	"github.com/imdonix/ncore-go/pkg/rest"
)

func main() {
	s := rest.NewServer()
	if err := s.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
