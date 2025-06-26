package postgres_test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/elokanugrah/go-order-system/internal/config"
	"github.com/elokanugrah/go-order-system/internal/database"
	"github.com/elokanugrah/go-order-system/internal/domain"
	"github.com/elokanugrah/go-order-system/internal/repository/postgres"
	"github.com/stretchr/testify/suite"
)

type ProductRepositorySuite struct {
	suite.Suite // Embed testify's suite package.

	db   *sql.DB
	repo *postgres.PostgresProductRepository
}

// SetupSuite runs once before all tests in this suite.
// It's used for setting up the database connection.
func (s *ProductRepositorySuite) SetupSuite() {
	// Load configuration to get database DSN
	// Note: Assumes .env file is in the project root.
	cfg := config.Load()

	// Create a new database connection for testing
	s.db = database.NewConnection(cfg)

	// Create the repository instance with the test database connection
	s.repo = postgres.NewProductRepository(s.db)
}

// TearDownSuite runs once after all tests in this suite are finished.
// It's used for cleaning up resources, like closing the database connection.
func (s *ProductRepositorySuite) TearDownSuite() {
	err := s.db.Close()
	if err != nil {
		log.Fatalf("Failed to close test database connection: %v", err)
	}
}

// TearDownTest runs after each test function in the suite.
// We use it to clean the database tables to ensure tests are isolated.
func (s *ProductRepositorySuite) TearDownTest() {
	// Truncate the table to leave it clean for the next test.
	_, err := s.db.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE")
	s.Suite.NoError(err)
}

// This function is the entry point for running the test suite.
func TestProductRepository(t *testing.T) {
	suite.Run(t, new(ProductRepositorySuite))
}

// TestSaveAndFindByID tests both Save and FindByID in a single flow.
func (s *ProductRepositorySuite) TestSaveAndFindByID() {
	// Use the suite's assertion library
	assert := s.Suite.Assert()
	ctx := context.Background()

	newProduct := &domain.Product{
		Name:     "Kopi Arabica",
		Price:    120000,
		Quantity: 50,
	}

	// Act 1: Save the new product to the database
	err := s.repo.Save(ctx, newProduct)

	// Assert 1: Check for errors and that the ID is now populated
	assert.NoError(err)
	assert.NotZero(newProduct.ID) // The ID should be populated by the DB

	// Act 2: Find the product we just saved using its new ID
	foundProduct, err := s.repo.FindByID(ctx, newProduct.ID)

	// Assert 2: Check for errors and that the found data matches
	assert.NoError(err)
	assert.NotNil(foundProduct)
	assert.Equal("Kopi Arabica", foundProduct.Name)
	assert.Equal(50, foundProduct.Quantity)
}

// TestFindByID_NotFound tests the case where a product ID does not exist.
func (s *ProductRepositorySuite) TestFindByID_NotFound() {
	assert := s.Suite.Assert()
	ctx := context.Background()

	// Act: Try to find a product with an ID that doesn't exist
	product, err := s.repo.FindByID(ctx, 99999)

	// Assert: We expect no SQL error, but the returned product should be nil
	assert.NoError(err)
	assert.Nil(product)
}

// TestUpdate tests if a product's data can be successfully updated.
func (s *ProductRepositorySuite) TestUpdate() {
	assert := s.Suite.Assert()
	ctx := context.Background()

	// Arrange: First, create a product to update
	productToUpdate := &domain.Product{Name: "Buku Lama", Price: 50000, Quantity: 5}
	err := s.repo.Save(ctx, productToUpdate)
	assert.NoError(err)

	// Act: Modify the product's details and call Update
	productToUpdate.Name = "Buku Baru Edisi Revisi"
	productToUpdate.Quantity = 3
	err = s.repo.Update(ctx, productToUpdate)
	assert.NoError(err)

	// Assert: Fetch the product again and verify its fields are updated
	updatedProduct, err := s.repo.FindByID(ctx, productToUpdate.ID)
	assert.NoError(err)
	assert.NotNil(updatedProduct)
	assert.Equal("Buku Baru Edisi Revisi", updatedProduct.Name)
	assert.Equal(3, updatedProduct.Quantity)
}

// TestDelete tests if a product can be successfully removed.
func (s *ProductRepositorySuite) TestDelete() {
	assert := s.Suite.Assert()
	ctx := context.Background()

	// Arrange: Create a product to delete
	productToDelete := &domain.Product{Name: "Barang Hapus", Price: 10, Quantity: 1}
	err := s.repo.Save(ctx, productToDelete)
	assert.NoError(err)

	// Act: Delete the product
	err = s.repo.Delete(ctx, productToDelete.ID)
	assert.NoError(err)

	// Assert: Try to find the deleted product, it should be nil
	foundProduct, err := s.repo.FindByID(ctx, productToDelete.ID)
	assert.NoError(err)
	assert.Nil(foundProduct)
}
