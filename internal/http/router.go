package http

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gingray/go-template/internal/http/handler"
	"github.com/gingray/go-template/internal/repo/postgres"
	"github.com/gingray/go-template/pkg/kafka"
)

type Router struct {
	userRepo      postgres.UserRepo
	kafkaProducer kafka.Client
}

func NewRouter(userRepo postgres.UserRepo, kafkaProducer kafka.Client) *Router {
	return &Router{userRepo: userRepo, kafkaProducer: kafkaProducer}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	probeHandler := handler.NewProbeHandler()
	router.GET("/health", probeHandler.HealthHandler)
	router.GET("/ready", probeHandler.ReadyHandler)
	router.GET("/", r.Root)
	router.GET("/produce", r.KafkaProducer)
}

func (r *Router) Root(c *gin.Context) {
	users, _ := r.userRepo.GetUsers(c)
	c.JSON(200, users)
}

func (r *Router) KafkaProducer(c *gin.Context) {
	users, _ := r.userRepo.GetUsers(c)
	r.kafkaProducer.Produce(context.Background(), "test-topic", "user_key:1", users)
	c.JSON(200, users)
}
