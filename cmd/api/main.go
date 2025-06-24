package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/elokanugrah/go-order-system/internal/config"
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

	// 3. Use the configuration to start the server
	// (Dependency Injection code would go here)

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil { // Replace nil with your router
		log.Fatalf("Failed to start server: %v", err)
	}
}
