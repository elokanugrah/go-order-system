package domain

import (
	"errors"
	"time"
)

var ErrInsufficientStock = errors.New("insufficient product stock")

type Product struct {
	ID        int64
	Name      string
	Price     float64
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsStockAvailable checks if the current stock is sufficient for the requested quantity.
func (p *Product) IsStockAvailable(requestedQuantity int) bool {
	return p.Quantity >= requestedQuantity
}

// DecreaseStock reduces the product's stock quantity.
func (p *Product) DecreaseStock(amount int) error {
	if !p.IsStockAvailable(amount) {
		return ErrInsufficientStock
	}
	if amount <= 0 {
		return errors.New("amount to decrease must be positive")
	}
	p.Quantity -= amount
	p.UpdatedAt = time.Now()
	return nil
}

// IncreaseStock increases the product's stock quantity.
func (p *Product) IncreaseStock(amount int) error {
	if amount <= 0 {
		return errors.New("amount to increase must be positive")
	}
	p.Quantity += amount
	p.UpdatedAt = time.Now()
	return nil
}
