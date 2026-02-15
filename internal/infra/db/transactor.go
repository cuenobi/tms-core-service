package db

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type contextKey string

const txKey contextKey = "db_tx"

// Transactor provides transaction management
type Transactor struct {
	db *gorm.DB
}

// NewTransactor creates a new Transactor
func NewTransactor(db *gorm.DB) *Transactor {
	return &Transactor{db: db}
}

// FromContext returns a database instance from context if it exists, otherwise returns defaultDB
func FromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}

// WithTransaction executes a function within a database transaction
func (t *Transactor) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Inject transaction into context
		txCtx := context.WithValue(ctx, txKey, tx)
		if err := fn(txCtx); err != nil {
			return fmt.Errorf("transaction failed: %w", err)
		}
		return nil
	})
}

// GetDB returns the underlying DB instance
func (t *Transactor) GetDB() *gorm.DB {
	return t.db
}
