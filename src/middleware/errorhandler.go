package middleware

import (
	"portfolio_golang/src/zaplog"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	// to capture all non-fatal/fatal error stack traces passed to c.Errors
	// to log meta key values only
	for _, err := range c.Errors {
		meta, ok := err.Meta.(map[string]string)
		if ok {
			fields := []zapcore.Field{
				zap.String("meta", meta["META"]),
				zap.String("handler", meta["HANDLER"]),
				zap.String("errortype", meta["ERRORTYPE"]),
			}
			zaplog.Logger.Error(err.Error(), fields...)
		}
		// other logic such as collating and sending to DB can be added
	}
}
