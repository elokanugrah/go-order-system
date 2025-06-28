package dto

type CreateProductInput struct {
	Name     string
	Price    float64
	Quantity int
}

type UpdateProductInput struct {
	Name     string
	Price    float64
	Quantity int
}
