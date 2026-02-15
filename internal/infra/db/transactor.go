package db

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Transactor provides transaction management
type Transactor struct {
	db *gorm.DB
}

// NewTransactor creates a new Transactor
func NewTransactor(db *gorm.DB) *Transactor {
	return &Transactor{db: db}
}

// WithTransaction executes a function within a database transaction
func (t *Transactor) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := fn(tx); err != nil {
			return fmt.Errorf("transaction failed: %w", err)
		}
		return nil
	})
}

// GetDB returns the underlying DB instance
func (t *Transactor) GetDB() *gorm.DB {
	return t.db
}
