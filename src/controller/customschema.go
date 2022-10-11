package controller

import (
	"net/http"
	"portfolio_golang/src/internal"

	"github.com/gin-gonic/gin"
)

type CustomSchemaPayload struct {
	VALIDATESTRING string   `binding:"required,alpha"`
	VALIDATESLICES []string `binding:"required"`
	VALIDATECUSTOM string   `binding:"required,testcustom"`
}

func CustomSchema(c *gin.Context) {
	data, ok := internal.StoredContextValues(c, "payload")
	if !ok {
		// abortWithStatusJSON is already set
		return
	}
	res := JSONResponse{
		Response: data["payload"],
		Message:  "schema is validated",
	}

	c.JSON(http.StatusOK, res)
}
