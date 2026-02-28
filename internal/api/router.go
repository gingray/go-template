package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gingray/go-template/pkg/config"
)

type Router struct {
	cfg *config.Config
}

func NewRouter(cfg *config.Config) *Router {
	return &Router{cfg: cfg}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	router.GET("/health", r.Ping)
	router.GET("/ready", r.Ready)
	router.GET("/", r.Root)
}

func (r *Router) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (r *Router) Ready(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ready",
	})
}

func (r *Router) Root(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world",
	})
}
