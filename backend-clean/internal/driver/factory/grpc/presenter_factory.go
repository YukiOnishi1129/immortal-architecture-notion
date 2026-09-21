// Package grpc provides factory functions for gRPC adapters.
package grpc

import grpcpresenter "immortal-architecture-cqrs/backend/internal/adapter/grpc/presenter"

// NewAccountOutputFactory returns a factory for gRPC AccountPresenter.
func NewAccountOutputFactory() func() *grpcpresenter.AccountPresenter {
	return func() *grpcpresenter.AccountPresenter {
		return grpcpresenter.NewAccountPresenter()
	}
}
