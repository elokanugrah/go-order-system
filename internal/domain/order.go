package domain

import (
	"errors"
	"time"
)

var ErrEmptyOrder = errors.New("order must have at least one item")

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusShipped   OrderStatus = "shipped"
	StatusCompleted OrderStatus = "completed"
	StatusCancelled OrderStatus = "cancelled"
)

// Order represents the core business entity for a customer's order.
type Order struct {
	ID          int64
	UserID      int64
	OrderItems  []OrderItem
	TotalAmount float64
	Status      OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// OrderItem represents a single line item within an order.
type OrderItem struct {
	ID           int64
	OrderID      int64
	Product      Product
	Quantity     int
	PriceAtOrder float64
}

// NewOrder is a constructor function to create a new Order.
// It ensures that every new order has sensible defaults and passes initial validation.
func NewOrder(userID int64, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}

	now := time.Now()

	// Create the order with default values.
	order := &Order{
		UserID:     userID,
		OrderItems: items,
		Status:     StatusPending, // Default status for any new order.
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Calculate the total amount upon creation.
	order.CalculateTotalAmount()

	return order, nil
}

// CalculateTotalAmount sums up the price of all items in the order.
func (o *Order) CalculateTotalAmount() {
	var total float64
	for _, item := range o.OrderItems {
		total += item.PriceAtOrder * float64(item.Quantity)
	}
	o.TotalAmount = total
}

// AddItem adds a new OrderItem to the order and recalculates the total amount.
func (o *Order) AddItem(item OrderItem) {
	o.OrderItems = append(o.OrderItems, item)
	o.CalculateTotalAmount() // Recalculate total after adding an item.
	o.UpdatedAt = time.Now()
}

// ChangeStatus updates the status of the order.
func (o *Order) ChangeStatus(newStatus OrderStatus) {
	o.Status = newStatus
	o.UpdatedAt = time.Now()
}
