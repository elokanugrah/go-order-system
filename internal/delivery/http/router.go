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
			products.GET("/:id", h.GetProductByID)
		}
	}

	return router
}
