package internal

import (
	"fmt"
	e "portfolio_golang/src/errors"

	"github.com/gin-gonic/gin"
)

// c.Get returns False if the key doesn't exist.
// To throw error instead.
// Takes in a variable number of string arg as keys.
func StoredContextValues(c *gin.Context, keys ...string) (map[string]interface{}, bool) {
	values := make(map[string]interface{}, len(keys))
	for _, k := range keys {
		value, ok := c.Get(k)
		if !ok {
			err := e.NewError(
				e.ErrMissingKey,
				fmt.Sprintf("%s key missing in gin.Context", k),
				c.HandlerName(),
			)
			c.Error(err).SetMeta(err.Meta())
			c.AbortWithStatusJSON(e.ErrorResponse(err))
			return nil, false
		}
		values[k] = value
	}
	return values, true
}
