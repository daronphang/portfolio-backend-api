package middleware

import (
	"github.com/gin-gonic/gin"
)

func RecoverAbortedRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// log.Printf("recovered from panic: client canceled the request: %s", c.Request.URL)
				// log.Println(err)
				c.Request.Context().Done()
			}
		}()
		c.Next()
	}
}
