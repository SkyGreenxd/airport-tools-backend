package config

import (
	"airport-tools-backend/pkg/e"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	const op = "LoadEnv"

	if err := godotenv.Load(); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}
