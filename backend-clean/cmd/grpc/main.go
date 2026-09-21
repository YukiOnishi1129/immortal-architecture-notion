// Package main starts the gRPC server.
package main

import (
	"context"
	"fmt"
	"log"

	initializer "immortal-architecture-cqrs/backend/internal/driver/initializer/grpc"
)

func main() {
	ctx := context.Background()
	s, cfg, cleanup, err := initializer.BuildServer(ctx)
	if err != nil {
		log.Fatalf("failed to initialize gRPC server: %v", err)
	}
	defer cleanup()

	lis, err := initializer.GetListener(cfg)
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	grpcPort := cfg.ServerPort + 1
	addr := fmt.Sprintf(":%d", grpcPort)
	log.Printf("starting gRPC server at %s\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
