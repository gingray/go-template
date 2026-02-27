package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gingray/go-template/pkg/common"
	"github.com/gingray/go-template/pkg/component"
)

type Server struct {
	component.BaseComponent
	router *gin.Engine
	logger *slog.Logger
	addr   string
	ch     chan struct{}
}

func (s *Server) Name() string {
	return "http_server"
}

func NewServer(cfg *common.HTTPServiceConfig, app *common.App) *Server {
	server := &Server{
		router: app.HttpRouter,
		logger: app.Logger,
		addr:   fmt.Sprintf(":%d", cfg.Port),
	}
	server.router.GET("/ping", func(c *gin.Context) {
		go func() {
			<-time.After(time.Second * 2)
			close(server.ch)
		}()
		c.JSON(200, gin.H{"message": "pong"})
	})

	server.AddReadyHandler(func(ctx context.Context) error {
		server.ch = make(chan struct{})
		return nil
	})
	return server
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting server")
	errCh := make(chan error)
	go func() {
		s.logger.Info("Server started", "addr", s.addr)

		srv := &http.Server{
			Addr:    s.addr,
			Handler: s.router.Handler(),
		}

		s.AddShutdownHandler(func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		})

		err := srv.ListenAndServe()
		errCh <- err
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	case <-s.ch:
		return fmt.Errorf("test shutdown")
	}
}
