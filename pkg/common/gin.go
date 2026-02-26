package common

import "github.com/gin-gonic/gin"

type HTTPServiceConfig struct {
	Port int `env:"HTTP_PORT" envDefault:"3000"`
}

func NewHTTPRouter() *gin.Engine {
	return gin.Default()
}
