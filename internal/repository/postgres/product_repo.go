package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/elokanugrah/go-order-system/internal/domain"
	"github.com/elokanugrah/go-order-system/internal/usecase"
)

// Ensure PostgresProductRepository implements the usecase.ProductRepository interface.
var _ usecase.ProductRepository = (*PostgresProductRepository)(nil)

// PostgresProductRepository is the PostgreSQL implementation of the ProductRepository interface.
type PostgresProductRepository struct {
	db *sql.DB
}

// FindByID retrieves a single product from the database by its ID.
func (r *PostgresProductRepository) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `SELECT id, name, price, quantity, created_at, updated_at FROM products WHERE id = $1`

	// Create a product variable to scan into.
	var p domain.Product

	// Execute the query. QueryRowContext is used because we expect exactly one row.
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Quantity,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	// Handle potential errors.
	if err != nil {
		// If no rows were found, it's a specific error we should handle.
		if errors.Is(err, sql.ErrNoRows) {
			// We can return a custom domain-specific error or nil, nil to indicate not found.
			// Returning nil, nil is simple and the use case can interpret it as "not found".
			return nil, nil
		}
		// For any other errors, wrap them with more context.
		return nil, fmt.Errorf("error scanning product: %w", err)
	}

	return &p, nil
}

// FindManyByIDs retrieves multiple products from the database by their IDs.
func (p *PostgresProductRepository) FindManyByIDs(ctx context.Context, ids []int64) ([]domain.Product, error) {
	// For now, we can leave it unimplemented as it's not needed for the "Get By ID" feature.
	return nil, errors.New("FindManyByIDs not implemented")
}

// Update updates a product's details in the database.
func (p *PostgresProductRepository) Update(ctx context.Context, product *domain.Product) error {
	// Implementation for this method would go here.
	return errors.New("Update not implemented")
}

// NewProductRepository creates a new instance of PostgresProductRepository.
func NewProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}
