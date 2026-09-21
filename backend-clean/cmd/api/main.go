// Package main starts the API server.
package main

import (
	"context"
	"log"
	"strconv"

	initializer "immortal-architecture-cqrs/backend/internal/driver/initializer/api"
)

func main() {
	ctx := context.Background()
	e, cfg, cleanup, err := initializer.BuildServer(ctx)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}
	defer cleanup()

	addr := ":" + strconv.Itoa(cfg.ServerPort)
	log.Printf("starting HTTP server at %s\n", addr)
	if err := e.Start(addr); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
