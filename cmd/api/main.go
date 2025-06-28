package main

import (
	"log"

	"github.com/elokanugrah/go-order-system/internal/config"
	"github.com/elokanugrah/go-order-system/internal/database"
	"github.com/elokanugrah/go-order-system/internal/messagebroker"
	"github.com/elokanugrah/go-order-system/internal/repository/postgres"
	"github.com/elokanugrah/go-order-system/internal/usecase"

	httpDelivery "github.com/elokanugrah/go-order-system/internal/delivery/http"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()

	db := database.NewConnection(cfg)
	defer db.Close()
	// --- WIRING / DEPENDENCY INJECTION ---

	// Initialize Message Broker
	mb, err := messagebroker.NewRabbitMQBroker(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize message broker: %v", err)
	}

	// Initialize Repository Layer
	productRepo := postgres.NewProductRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	txManager := postgres.NewTransactionManager(db)

	// Initialize Usecase Layer
	productUseCase := usecase.NewProductUseCase(productRepo)
	orderUseCase := usecase.NewOrderUseCase(orderRepo, productRepo, txManager, mb)

	// Initialize Delivery Layer (Handler)
	// For now, orderUseCase is nil because we haven't built it completely.
	apiHandler := httpDelivery.NewHandler(productUseCase, orderUseCase)

	// Setup Router and Start Server
	router := httpDelivery.SetupRouter(apiHandler)

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
