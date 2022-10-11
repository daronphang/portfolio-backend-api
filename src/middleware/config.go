package middleware

import (
	"portfolio_golang/src/config"

	"github.com/gin-gonic/gin"
)

func SetConfig(appCfg *config.EnvConf) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("appCfg", appCfg)
		c.Next()
	}
}
