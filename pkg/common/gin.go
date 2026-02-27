package common

import (
	"github.com/gin-gonic/gin"
)

type HTTPServiceConfig struct {
	Port int `env:"HTTP_PORT" envDefault:"3000"`
}

func (a *App) WithHTTPRouter() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		a.Logger.Info("GIN request", "method", param.Method, "path", param.Path, "code", param.StatusCode,
			"latency", param.Latency.String(), "ip", param.ClientIP, "error", param.ErrorMessage)
		return ""
	}))
	a.HttpRouter = router
	return nil
}
