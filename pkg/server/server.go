package server

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gingray/go-template/pkg/common"
	"github.com/gingray/go-template/pkg/component"
)

type Server struct {
	component.BaseComponent
	router *gin.Engine
	logger *slog.Logger
}

func NewServer(app *common.App) *Server {
	return &Server{
		router: app.HttpRouter,
		logger: app.Logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting server")
	errCh := make(chan error)
	go func() {
		s.router.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		err := s.router.Run()
		errCh <- err
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
