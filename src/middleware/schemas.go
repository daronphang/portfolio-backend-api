package middleware

import (
	"errors"
	"fmt"

	"portfolio_golang/src/controller"
	e "portfolio_golang/src/errors"
	"portfolio_golang/src/internal"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateSchema(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload interface{}
		var schema string

		if len(name) == 0 {
			// extract URL param
			schema = c.Param("name")
		} else {
			schema = name
		}

		switch strings.ToUpper(schema) {
		case "CUSTOMSCHEMA":
			payload = &controller.CustomSchemaPayload{}
		case "CREATECONTACT":
			payload = &controller.ContactPayload{}
		default:
			err := e.NewError(
				e.ErrMissingKey,
				fmt.Sprintf("%s schema provided does not exist", schema),
				c.HandlerName(),
			)
			c.Error(err).SetMeta(err.Meta())
			c.AbortWithStatusJSON(e.ErrorResponse(err))
			return
		}

		if err := c.ShouldBindJSON(payload); err != nil {
			// gin uses playground/validator for validation
			// returns an array of type FieldError
			var ve validator.ValidationErrors
			var error string
			if !errors.As(err, &ve) {
				error = err.Error()
			} else {
				error = internal.MapToString(e.TranslateErrors(&ve))
			}

			err := e.NewError(e.ErrInvalidSchema, error, c.HandlerName())
			c.Error(err).SetMeta(err.Meta())
			c.AbortWithStatusJSON(e.ErrorResponse(err))
			return
		}

		// payload is valid and parsed, passed to next middleware
		c.Set("payload", payload)
		c.Next()
	}
}
