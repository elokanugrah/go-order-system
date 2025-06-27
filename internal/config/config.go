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

// The `env` tag is used to specify the environment variable name.
type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`

	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
}

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func Load() *Config {
	if err := loadDotEnv(); err != nil {
		if !os.IsNotExist(err) {
			log.Println("Error loading .env file, but it's not a 'file not found' error:", err)
		}
	}

	cfg := Config{}
	// The `env.Parse` function will use the `env` tags to map variables.
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse configuration: %+v", err)
	}

	return &cfg
}

// loadDotEnv searches for a .env file from the current directory up to the root
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
