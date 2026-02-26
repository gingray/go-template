package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gingray/go-template/pkg/common"
	"github.com/gingray/go-template/pkg/component"
)

type Server struct {
	component.BaseComponent
	router *gin.Engine
	logger *slog.Logger
	addr   string
}

func (s *Server) Name() string {
	return "http_server"
}

func NewServer(cfg *common.HTTPServiceConfig, app *common.App) *Server {
	return &Server{
		router: app.HttpRouter,
		logger: app.Logger,
		addr:   fmt.Sprintf(":%d", cfg.Port),
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting server")
	errCh := make(chan error)
	go func() {
		s.router.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		s.logger.Info("Server started", "addr", s.addr)
		err := s.router.Run(s.addr)
		errCh <- err
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
