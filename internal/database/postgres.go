package database

import (
	"database/sql"
	"log"

	"github.com/elokanugrah/go-order-system/internal/config"
)

// NewConnection creates and returns a new database connection pool (*sql.DB).
func NewConnection(cfg *config.Config) *sql.DB {
	// Use the DSN method from the config struct to get the connection string.
	dsn := cfg.DSN()

	// sql.Open just validates its arguments, it doesn't create a connection yet.
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("FATAL: Could not prepare database connection: %v", err)
	}

	// Ping the database to verify that a connection is established and the credentials are valid.
	if err := db.Ping(); err != nil {
		log.Fatalf("FATAL: Database is not reachable: %v", err)
	}

	log.Println("Database connection successful.")
	return db
}
