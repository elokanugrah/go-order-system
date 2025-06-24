package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/elokanugrah/go-order-system/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	productUseCase *usecase.ProductUseCase
	orderUseCase   *usecase.OrderUseCase
}

func NewHandler(puc *usecase.ProductUseCase, ouc *usecase.OrderUseCase) *Handler {
	return &Handler{
		productUseCase: puc,
		orderUseCase:   ouc,
	}
}

// GetProductByID handles the HTTP request for fetching a single product.
func (h *Handler) GetProductByID(c *gin.Context) {
	// 1. Get the ID from the URL parameter.
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	// 2. Call the use case to get the product.
	// We use c.Request.Context() to pass the request context down through the layers.
	product, err := h.productUseCase.GetProductByID(c.Request.Context(), id)
	if err != nil {
		// 3. Handle errors returned from the use case.
		if errors.Is(err, usecase.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		// For any other unexpected errors, return a 500.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	// 4. On success, return the product data with a 200 OK status.
	// We can define a DTO here to control the output format if needed,
	// but for simplicity, we'll return the domain object directly.
	c.JSON(http.StatusOK, product)
}
