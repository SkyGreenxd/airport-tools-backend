package server

import (
	"airport-tools-backend/internal/config"
	"context"
	"net/http"

	"github.com/rs/cors"
)

// Server обёртка над http.Server для запуска и остановки HTTP-сервиса.
type Server struct {
	httpServer *http.Server
}

// NewServer создаёт новый HTTP-сервер с заданным обработчиком и конфигурацией.
func NewServer(handler http.Handler, httpServer config.HttpServer) *Server {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: true,
	}).Handler(handler)

	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + httpServer.Port,
			Handler:      corsHandler,
			ReadTimeout:  httpServer.ReadTimeout,
			WriteTimeout: httpServer.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
