package rest

import (
	"portfolio_golang/src/controller"
	"portfolio_golang/src/middleware"

	"github.com/gin-gonic/gin"
)

func addTesting(rg *gin.RouterGroup) {
	testing := rg.Group("/testing")

	{
		testing.GET(
			"/heartbeat",
			controller.HeartBeat,
		)

		testing.POST(
			"/customschema",
			middleware.ValidateSchema("CUSTOMSCHEMA"),
			controller.CustomSchema,
		)
	}
}