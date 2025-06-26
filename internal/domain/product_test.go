package domain_test

import (
	"testing"

	"github.com/elokanugrah/go-order-system/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestProduct_DecreaseStock(t *testing.T) {

	t.Run("should decrease stock successfully when stock is sufficient", func(t *testing.T) {
		// Arrange: Set up the initial state
		product := &domain.Product{
			ID:       1,
			Name:     "Test Book",
			Quantity: 10,
		}
		amountToDecrease := 4

		// Act: Perform the action to be tested
		err := product.DecreaseStock(amountToDecrease)

		// Assert: Check if the outcome is as expected
		assert.NoError(t, err)               // We expect no error
		assert.Equal(t, 6, product.Quantity) // We expect the quantity to be 10 - 4 = 6
	})

	t.Run("should return an error when stock is insufficient", func(t *testing.T) {
		// Arrange: Set up the initial state
		product := &domain.Product{
			ID:       2,
			Name:     "Test Gadget",
			Quantity: 5,
		}
		amountToDecrease := 7 // Trying to decrease more than available

		// Act: Perform the action
		err := product.DecreaseStock(amountToDecrease)

		// Assert: Check for the expected error and state
		assert.Error(t, err)                              // We expect an error
		assert.Equal(t, domain.ErrInsufficientStock, err) // We expect the specific error
		assert.Equal(t, 5, product.Quantity)              // The quantity should not change
	})

	t.Run("should return an error for non-positive decrease amount", func(t *testing.T) {
		// Arrange
		product := &domain.Product{Quantity: 10}

		// Act
		err := product.DecreaseStock(-1)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, 10, product.Quantity)
	})
}

func TestProduct_IncreaseStock(t *testing.T) {

	t.Run("should increase stock successfully with a positive amount", func(t *testing.T) {
		// Arrange: Set up the initial state of the product.
		product := &domain.Product{
			ID:       1,
			Name:     "Test Monitor",
			Quantity: 20,
		}
		amountToIncrease := 10

		// Act: Perform the action to be tested.
		err := product.IncreaseStock(amountToIncrease)

		// Assert: Verify the outcome.
		assert.NoError(t, err)                // We expect no error.
		assert.Equal(t, 30, product.Quantity) // We expect the quantity to be 20 + 10 = 30.
	})

	t.Run("should return an error when increasing by zero", func(t *testing.T) {
		// Arrange
		product := &domain.Product{
			ID:       2,
			Name:     "Test Keyboard",
			Quantity: 15,
		}
		amountToIncrease := 0

		// Act
		err := product.IncreaseStock(amountToIncrease)

		// Assert
		assert.Error(t, err)                  // We expect an error.
		assert.Equal(t, 15, product.Quantity) // The quantity should remain unchanged.
		assert.Equal(t, "amount to increase must be positive", err.Error())
	})

	t.Run("should return an error when increasing by a negative amount", func(t *testing.T) {
		// Arrange
		product := &domain.Product{
			ID:       3,
			Name:     "Test Mouse",
			Quantity: 50,
		}
		amountToIncrease := -5

		// Act
		err := product.IncreaseStock(amountToIncrease)

		// Assert
		assert.Error(t, err)                  // We expect an error.
		assert.Equal(t, 50, product.Quantity) // The quantity must not change.
		assert.Equal(t, "amount to increase must be positive", err.Error())
	})
}
