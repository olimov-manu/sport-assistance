package server

import (
	"context"
	"net/http"
	"sport-assistance/pkg/configs"

	"github.com/bytedance/gopkg/util/logger"
)

type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func NewServer(handler http.Handler, cfg *configs.Config) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.ServerConfig.Port,
			Handler:        handler,
			ReadTimeout:    cfg.ServerConfig.ReadTimeout,
			WriteTimeout:   cfg.ServerConfig.WriteTimeout,
			MaxHeaderBytes: 1 << 20, // 1 MB
		},
	}
}
