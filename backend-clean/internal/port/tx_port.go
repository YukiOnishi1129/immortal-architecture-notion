package port

import "context"

// TxManager provides transaction boundary control.
type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
