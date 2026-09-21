// Package factory provides constructors for driver-level wiring.
package factory

import "immortal-architecture-notion/backend/internal/port"

// NewTxFactory returns a factory that provides TxManager.
func NewTxFactory(tx port.TxManager) func() port.TxManager {
	return func() port.TxManager {
		return tx
	}
}
