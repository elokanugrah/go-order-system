package usecase

import (
	"context"
	"errors"

	"github.com/elokanugrah/go-order-system/internal/domain"
	"github.com/elokanugrah/go-order-system/internal/dto"
)

// OrderUseCase handles the business logic for orders.
// It depends on repository interfaces to interact with the data layer.
type OrderUseCase struct {
	orderRepo   OrderRepository
	productRepo ProductRepository
	txManager   TransactionManager // Dependency for transaction management
}

// NewOrderUseCase is the constructor for OrderUseCase.
// It receives dependencies (repositories, etc.) and returns a new OrderUseCase instance.
func NewOrderUseCase(or OrderRepository, pr ProductRepository, tm TransactionManager) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:   or,
		productRepo: pr,
		txManager:   tm,
	}
}

// CreateOrder is the primary method for creating a new order.
// It orchestrates fetching products, validating stock, creating domain objects, and persisting them.
func (uc *OrderUseCase) CreateOrder(ctx context.Context, input dto.CreateOrderInput) (*domain.Order, error) {
	// --- Input Validation ---
	if len(input.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	var createdOrder *domain.Order

	// --- Transactional Business Logic ---
	// Execute the entire creation process within a single database transaction
	// using the callback pattern provided by our TransactionManager.
	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Get all product IDs from the input to fetch them in one query.
		productIDs := make([]int64, len(input.Items))
		itemMap := make(map[int64]dto.CreateOrderItemInput)
		for i, item := range input.Items {
			if item.Quantity <= 0 {
				return errors.New("item quantity must be positive")
			}
			productIDs[i] = item.ProductID
			itemMap[item.ProductID] = item
		}

		// 2. Fetch all required products from the database.
		// We pass the transactional context `txCtx` here.
		products, err := uc.productRepo.FindManyByIDs(txCtx, productIDs)
		if err != nil {
			return err
		}
		if len(products) != len(productIDs) {
			return errors.New("one or more products not found")
		}

		var orderItems []domain.OrderItem
		var productsToUpdate []*domain.Product

		// 3. Validate stock and prepare domain objects.
		for _, p := range products {
			itemInput := itemMap[p.ID]

			// Use the business logic from the domain entity to check stock.
			if !p.IsStockAvailable(itemInput.Quantity) {
				return domain.ErrInsufficientStock // Use the domain-specific error.
			}

			// Decrease the stock using the domain entity method.
			if err := p.DecreaseStock(itemInput.Quantity); err != nil {
				return err
			}

			orderItems = append(orderItems, domain.OrderItem{
				Product:      p,
				Quantity:     itemInput.Quantity,
				PriceAtOrder: p.Price,
			})

			// Add the product to a list to be updated later.
			productToUpdate := p
			productsToUpdate = append(productsToUpdate, &productToUpdate)
		}

		// 4. Create the main Order domain object.
		createdOrder, err = domain.NewOrder(input.UserID, orderItems)
		if err != nil {
			return err
		}

		// 5. Persist the order and its items to the database.
		// We pass the transactional context `txCtx` here.
		if err := uc.orderRepo.Save(txCtx, createdOrder); err != nil {
			return err
		}

		// 6. Persist the updated product stock for all affected products.
		// We pass the transactional context `txCtx` here.
		for _, p := range productsToUpdate {
			if err := uc.productRepo.Update(txCtx, p); err != nil {
				return err
			}
		}

		// If success, the TransactionManager will commit the transaction.
		return nil
	})

	// If the transaction failed return err
	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}
