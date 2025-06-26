package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
// It is populated from environment variables and/or a .env file.
// The `env` tag is used to specify the environment variable name.
type Config struct {
	// Server configuration
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`

	// Database configuration
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
}

// DSN generates the Data Source Name string for connecting to the PostgreSQL database.
func (c *Config) DSN() string {
	// Example: "postgres://user:password@localhost:5432/order_db?sslmode=disable"
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func Load() *Config {
	// Attempt to load .env file.
	// This is useful for local development. In a production environment,
	// environment variables should be set directly.
	// We ignore the error if the file doesn't exist.
	if err := loadDotEnv(); err != nil {
		// Check if the error is other than the file not existing
		if !os.IsNotExist(err) {
			log.Println("Error loading .env file, but it's not a 'file not found' error:", err)
		}
	}

	cfg := Config{}
	// Parse environment variables into the Config struct.
	// The `env.Parse` function will use the `env` tags to map variables.
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse configuration: %+v", err)
	}

	return &cfg
}

// loadDotEnv searches for a .env file from the current directory up to the root
// and loads it if found.
func loadDotEnv() error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Search up to 5 levels up for the .env file.
	for i := 0; i < 5; i++ {
		envPath := filepath.Join(cwd, ".env")
		if _, err := os.Stat(envPath); err == nil {
			// .env file found, load it.
			return godotenv.Load(envPath)
		}
		// Go one directory up.
		cwd = filepath.Dir(cwd)
	}

	return errors.New(".env file not found in parent directories")
}
