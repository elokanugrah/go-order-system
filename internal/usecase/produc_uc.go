package usecase

import (
	"context"
	"errors"

	"github.com/elokanugrah/go-order-system/internal/domain"
)

var ErrProductNotFound = errors.New("product not found")

type ProductUseCase struct {
	productRepo ProductRepository
}

func NewProductUseCase(pr ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: pr,
	}
}

// GetProductByID contains the logic for retrieving a single product.
// For now, it directly calls the repository, but in the future,
// this is where caching, authorization, or other business rules would go.
func (uc *ProductUseCase) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	// Call the repository method to find the product.
	product, err := uc.productRepo.FindByID(ctx, id)
	if err != nil {
		// Handle repository-level errors.
		return nil, err
	}

	// If the repository returns nil, it means the product was not found.
	// We return a use-case-specific error.
	if product == nil {
		return nil, ErrProductNotFound
	}

	return product, nil
}
