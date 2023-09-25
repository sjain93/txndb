package config

import (
	"github.com/joho/godotenv"
)

func LoadEnvVars() error {
	return godotenv.Load(".env")
}
