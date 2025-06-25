package main

import (
	"database/sql"
	"log"

	"github.com/elokanugrah/go-order-system/internal/config"
	"github.com/elokanugrah/go-order-system/internal/repository/postgres"
	"github.com/elokanugrah/go-order-system/internal/usecase"

	httpDelivery "github.com/elokanugrah/go-order-system/internal/delivery/http"

	_ "github.com/lib/pq"
)

func main() {
	// 1. Load configuration
	cfg := config.Load()

	// 2. Use the configuration to connect to the database
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is not reachable: %v", err)
	}

	log.Println("Database connection successful.")

	// --- WIRING / DEPENDENCY INJECTION ---

	// 3. Initialize Repository Layer
	productRepo := postgres.NewProductRepository(db)

	// 4. Initialize Usecase Layer
	productUseCase := usecase.NewProductUseCase(productRepo)

	// 5. Initialize Delivery Layer (Handler)
	// For now, orderUseCase is nil because we haven't built it completely.
	apiHandler := httpDelivery.NewHandler(productUseCase, nil)

	// 6. Setup Router and Start Server
	router := httpDelivery.SetupRouter(apiHandler)

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
