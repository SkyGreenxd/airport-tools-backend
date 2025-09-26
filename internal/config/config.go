package config

import (
	"airport-tools-backend/pkg/e"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	const op = "LoadEnv"

	if err := godotenv.Load(); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

type HttpServer struct {
	Port         string        `mapstructure:"HTTP_PORT"`
	ReadTimeout  time.Duration `mapstructure:"HTTP_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"HTTP_WRITE_TIMEOUT"`
}

func LoadHttpServerConfig() HttpServer {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	readTimeout, _ := time.ParseDuration(os.Getenv("HTTP_READ_TIMEOUT"))
	writeTimeout, _ := time.ParseDuration(os.Getenv("HTTP_WRITE_TIMEOUT"))

	return HttpServer{
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}
