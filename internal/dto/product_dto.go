package dto

// --- Data Transfer Objects (DTOs) for Inputs ---

// CreateProductInput defines the data needed to create a new product.
type CreateProductInput struct {
	Name     string
	Price    float64
	Quantity int
}

// UpdateProductInput defines the data needed to update an existing product.
type UpdateProductInput struct {
	Name     string
	Price    float64
	Quantity int
}
