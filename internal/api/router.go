package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gingray/go-template/pkg/app"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Router struct {
	pgPool      *pgxpool.Pool
	kafkaClient *kgo.Client
}

func NewRouter(app *app.App) *Router {
	return &Router{pgPool: app.PgPool, kafkaClient: app.Kafka}
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
	value, _ := json.Marshal(msg)
	record := kgo.Record{
		Key:   []byte(key),
		Topic: "test-topic",
		Value: value,
	}
	r.kafkaClient.ProduceSync(context.Background(), &record)
	c.JSON(200, msg)
}
