package common

import "github.com/gin-gonic/gin"

type HTTPServiceConfig struct {
	Port int `env:"HTTP_PORT" envDefault:"3000"`
}

func (a *App) WithHTTPRouter() error {
	a.HttpRouter = gin.Default()
	return nil
}
