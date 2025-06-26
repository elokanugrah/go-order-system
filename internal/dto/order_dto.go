package dto

type CreateOrderItemInput struct {
	ProductID int64
	Quantity  int
}

type CreateOrderInput struct {
	UserID int64
	Items  []CreateOrderItemInput
}
