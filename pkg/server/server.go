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

func NewServer(handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + configs.GetConfigs().ServerConfig.Port,
			Handler:        handler,
			ReadTimeout:    configs.GetConfigs().ServerConfig.ReadTimeout,
			WriteTimeout:   configs.GetConfigs().ServerConfig.WriteTimeout,
			MaxHeaderBytes: 1 << 20, // 1 MB
		},
	}
}
