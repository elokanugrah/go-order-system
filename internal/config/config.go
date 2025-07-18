package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"9000"`

	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`

	RabbitMQURL string `env:"RABBITMQ_URL,required"`
}

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func Load() *Config {
	loadDotEnv()

	cfg := Config{}
	// The `env.Parse` function will use the `env` tags to map variables.
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse configuration: %+v", err)
	}

	return &cfg
}

// loadDotEnv searches for a .env file from the current directory up to the root
func loadDotEnv() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("WARN: Could not get working directory to find .env: %v", err)
		return
	}

	// Search up to 5 levels up for the .env file.
	for range 5 {
		envPath := filepath.Join(cwd, ".env")
		if _, err := os.Stat(envPath); err == nil {
			// .env file found, load it and then return.
			if err := godotenv.Load(envPath); err != nil {
				log.Printf("WARN: Error loading .env file from %s: %v", envPath, err)
			}
			return // Exit after finding and attempting to load.
		}
		// Go one directory up.
		cwd = filepath.Dir(cwd)
	}

	log.Printf("INFO: .env file not found. Relying on system-set environment variables.")
}
