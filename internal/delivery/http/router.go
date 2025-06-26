package http

import "github.com/gin-gonic/gin"

func SetupRouter(h *Handler) *gin.Engine {
	router := gin.Default()

	// Group routes under /api/v1
	api := router.Group("/api/v1")
	{
		// Group product-related routes
		products := api.Group("/products")
		{
			products.POST("/", h.CreateProduct)
			products.GET("/", h.ListProducts)
			products.GET("/:id", h.GetProductByID)
			products.PUT("/:id", h.UpdateProduct)
			products.DELETE("/:id", h.DeleteProduct)
		}
	}

	return router
}
