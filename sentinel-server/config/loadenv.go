package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(key string) (string, error) {
	// check if prod
	prod := os.Getenv("PROD")

	if prod != "true" {
		err := godotenv.Load()
		if err != nil {
			return "", err
		}
	}

	return os.Getenv(key), nil
}
