// Package initializer wires dependencies for the gRPC server.
package initializer

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	grpccontroller "immortal-architecture-cqrs/backend/internal/adapter/grpc/controller"
	"immortal-architecture-cqrs/backend/internal/adapter/grpc/generated/accountpb"
	"immortal-architecture-cqrs/backend/internal/driver/config"
	driverdb "immortal-architecture-cqrs/backend/internal/driver/db"
	"immortal-architecture-cqrs/backend/internal/driver/factory"
	grpcfactory "immortal-architecture-cqrs/backend/internal/driver/factory/grpc"
)

// BuildServer composes all dependencies and returns a gRPC server, config, and cleanup function.
func BuildServer(ctx context.Context) (*grpc.Server, *config.Config, func(), error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, func() {}, err
	}

	pool, err := driverdb.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, func() {}, err
	}
	cleanup := func() {
		pool.Close()
	}

	accountRepoFactory := factory.NewAccountRepoFactory(pool)
	accountInputFactory := factory.NewAccountInputFactory()
	accountOutputFactory := grpcfactory.NewAccountOutputFactory()

	// Create gRPC server
	s := grpc.NewServer()

	// Register account service
	accountController := grpccontroller.NewAccountController(
		accountInputFactory,
		accountOutputFactory,
		accountRepoFactory,
	)
	accountpb.RegisterAccountServiceServer(s, accountController)

	return s, cfg, cleanup, nil
}

// GetListener creates a TCP listener for the gRPC server.
func GetListener(cfg *config.Config) (net.Listener, error) {
	// gRPC uses a different port, add 1 to HTTP port for gRPC
	grpcPort := cfg.ServerPort + 1
	addr := fmt.Sprintf(":%d", grpcPort)
	return net.Listen("tcp", addr)
}
