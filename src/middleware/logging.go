package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"portfolio_golang/src/zaplog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// for retrieving error JSON responses
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinLogger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	body, _ := ioutil.ReadAll(c.Request.Body)
	// write body back to the request after reading it in middleware
	// returns an object with Close()
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

	// implementing custom Writer to intercept Write() calls
	blw := &bodyLogWriter{body: new(bytes.Buffer), ResponseWriter: c.Writer}
	c.Writer = blw

	// passing middleware
	c.Next()

	// logging logic
	var payload string
	cost := time.Since(start)

	// converting JSON request payload to string
	if len(body) > 0 {
		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, body); err != nil {
			zaplog.Logger.Warn("unable to compact JSON payload for logging")
		} else {
			payload = buffer.String()
		}
	}

	// [scheme:][//[userinfo@]host][/]path[?query][#fragment]

	fields := []zapcore.Field{
		zap.Int("status", c.Writer.Status()),
		zap.String("path", path),
		zap.String("method", c.Request.Method),
		zap.Duration("cost", cost),
		zap.String("clientip", c.ClientIP()),
		zap.String("queryparams", raw),
		zap.String("payload", strings.ReplaceAll(payload, "\"", "'")),
	}

	if c.Writer.Status() > 300 {
		buffer := new(bytes.Buffer)
		_ = json.Compact(buffer, blw.body.Bytes())
		fields = append(fields, zap.String("errbody", strings.ReplaceAll(buffer.String(), "\"", "'")))
		zaplog.Logger.Error("error logging", fields...)
	} else {
		zaplog.Logger.Info("after request logging", fields...)
	}
}
