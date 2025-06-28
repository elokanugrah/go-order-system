package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/elokanugrah/go-order-system/internal/domain"
	"github.com/elokanugrah/go-order-system/internal/dto"
	"github.com/elokanugrah/go-order-system/internal/usecase"
	"github.com/elokanugrah/go-order-system/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderUseCase_CreateOrder(t *testing.T) {
	var mockProductRepo *mocks.ProductRepository
	var mockOrderRepo *mocks.OrderRepository
	var mockTxManager *mocks.TransactionManager
	var orderUseCase *usecase.OrderUseCase

	// setup is a helper function to initialize components for each test.
	setup := func() {
		mockProductRepo = new(mocks.ProductRepository)
		mockOrderRepo = new(mocks.OrderRepository)
		mockTxManager = new(mocks.TransactionManager)
		orderUseCase = usecase.NewOrderUseCase(mockOrderRepo, mockProductRepo, mockTxManager)
	}

	t.Run("should create order successfully when all conditions are met", func(t *testing.T) {
		setup() // Reset mocks

		// Define the input from the delivery layer
		input := dto.CreateOrderInput{
			UserID: 123,
			Items: []dto.CreateOrderItemInput{
				{ProductID: 1, Quantity: 2},
				{ProductID: 2, Quantity: 1},
			},
		}

		// Define the data our mock ProductRepository should return
		mockProducts := []domain.Product{
			{ID: 1, Name: "Product A", Price: 10000, Quantity: 10}, // Stock is sufficient (10 > 2)
			{ID: 2, Name: "Product B", Price: 5000, Quantity: 5},   // Stock is sufficient (5 > 1)
		}

		// Arrange: Set up the TransactionManager mock
		// When WithTransaction is called, we immediately execute the function passed to it.
		mockTxManager.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
			Return(nil).
			Run(func(args mock.Arguments) {
				// Get the callback function and execute it
				callback := args.Get(1).(func(ctx context.Context) error)
				callback(context.Background())
			}).Once()

		// Arrange: Set up expectations for calls made *inside* the transaction
		mockProductRepo.On("FindManyByIDs", mock.Anything, []int64{1, 2}).Return(mockProducts, nil).Once()
		mockProductRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil).Times(2) // Expect Update to be called twice
		mockOrderRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil).Once()

		// Act
		createdOrder, err := orderUseCase.CreateOrder(context.Background(), input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, createdOrder)
		assert.Equal(t, float64(25000), createdOrder.TotalAmount) // (2 * 10000) + (1 * 5000)
		assert.Equal(t, domain.StatusPending, createdOrder.Status)

		// Assert that all mock expectations
		mockProductRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("should return error when stock is insufficient", func(t *testing.T) {
		setup()

		// Arrange: Input where quantity requested is more than stock
		input := dto.CreateOrderInput{UserID: 123, Items: []dto.CreateOrderItemInput{{ProductID: 1, Quantity: 20}}}
		mockProducts := []domain.Product{{ID: 1, Name: "Product A", Price: 10000, Quantity: 10}} // Only 10 in stock

		// Arrange: Transaction will still start
		mockTxManager.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
			Return(domain.ErrInsufficientStock). // We expect the whole transaction to fail with this specific error
			Run(func(args mock.Arguments) {
				callback := args.Get(1).(func(ctx context.Context) error)
				callback(context.Background())
			}).Once()

		// Arrange: Mock the call inside the transaction
		mockProductRepo.On("FindManyByIDs", mock.Anything, []int64{1}).Return(mockProducts, nil).Once()

		// Act
		createdOrder, err := orderUseCase.CreateOrder(context.Background(), input)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInsufficientStock) // Check for our specific domain error
		assert.Nil(t, createdOrder)

		// Assert that Save and Update were never called because logic failed early
		mockOrderRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
		mockProductRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("should return error when product is not found", func(t *testing.T) {
		setup()

		input := dto.CreateOrderInput{UserID: 123, Items: []dto.CreateOrderItemInput{{ProductID: 99, Quantity: 1}}}

		// Arrange: Mock FindManyByIDs to return an empty slice, simulating a not-found product
		mockProducts := []domain.Product{}
		mockTxManager.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
			Return(errors.New("one or more products not found")).
			Run(func(args mock.Arguments) {
				callback := args.Get(1).(func(ctx context.Context) error)
				callback(context.Background())
			}).Once()

		mockProductRepo.On("FindManyByIDs", mock.Anything, []int64{99}).Return(mockProducts, nil).Once()

		// Act
		createdOrder, err := orderUseCase.CreateOrder(context.Background(), input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, createdOrder)
		assert.Contains(t, err.Error(), "products not found")
	})
}
