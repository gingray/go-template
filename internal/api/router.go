package api

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gingray/go-template/pkg/app"
	"github.com/gingray/go-template/pkg/kafka"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct {
	pgPool        *pgxpool.Pool
	kafkaProducer kafka.Client
}

func NewRouter(app *app.App, kafkaProducer kafka.Client) *Router {
	return &Router{pgPool: app.PgPool, kafkaProducer: kafkaProducer}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	router.GET("/health", r.Ping)
	router.GET("/ready", r.Ready)
	router.GET("/", r.Root)
	router.GET("/produce", r.KafkaProducer)

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
	var email string
	var id int64
	r.pgPool.QueryRow(c, "SELECT id, email from users").Scan(&id, &email)
	c.JSON(200, gin.H{
		"id":    id,
		"email": email,
	})
}

func (r *Router) KafkaProducer(c *gin.Context) {
	var email string
	var id int64
	r.pgPool.QueryRow(c, "SELECT id, email from users").Scan(&id, &email)
	key := fmt.Sprintf("user_id:%d", id)
	msg := map[string]interface{}{
		"id":    id,
		"email": email,
	}
	r.kafkaProducer.Produce(context.Background(), "test-topic", key, msg)
	c.JSON(200, msg)
}
