package server

import (
	"airport-tools-backend/internal/config"
	"context"
	"net/http"
)

// Server обёртка над http.Server для запуска и остановки HTTP-сервиса.
type Server struct {
	httpServer *http.Server
}

// NewServer создаёт новый HTTP-сервер с заданным обработчиком и конфигурацией.
func NewServer(handler http.Handler, httpServer config.HttpServer) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + httpServer.Port,
			Handler:      handler,
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
