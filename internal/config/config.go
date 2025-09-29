package config

import (
	"os"
	"time"
)

const (
	defaultPort = "8080"
)

type HttpServer struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// LoadHttpServerConfig загружает конфигурацию HTTP-сервера из переменных окружения
func LoadHttpServerConfig() HttpServer {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = defaultPort
	}

	readTimeout, _ := time.ParseDuration(os.Getenv("HTTP_READ_TIMEOUT"))
	writeTimeout, _ := time.ParseDuration(os.Getenv("HTTP_WRITE_TIMEOUT"))

	return HttpServer{
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}
