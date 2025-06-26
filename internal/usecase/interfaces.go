package usecase

import (
	"context"

	"github.com/elokanugrah/go-order-system/internal/domain"
)

type ProductRepository interface {
	// Create
	Save(ctx context.Context, product *domain.Product) error

	// Read
	FindByID(ctx context.Context, id int64) (*domain.Product, error)
	FindManyByIDs(ctx context.Context, ids []int64) ([]domain.Product, error)
	FindAll(ctx context.Context, limit, offset int) ([]domain.Product, error)

	// Update
	Update(ctx context.Context, product *domain.Product) error

	// Delete
	Delete(ctx context.Context, id int64) error
}

type OrderRepository interface {
	// Create
	Save(ctx context.Context, order *domain.Order) error
}

// TransactionManager defines the contract for database transaction management.
// This allows use cases to run operations within a single transaction
// without being coupled to a specific database implementation.
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
