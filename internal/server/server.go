package server

import (
	"airport-tools-backend/internal/config"
	"context"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

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
