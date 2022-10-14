package rest

import (
	"portfolio_golang/src/customvalidator"
	"portfolio_golang/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RunGinWebApp() *gin.Engine {
	r := gin.New()
	r.Use(middleware.GinLogger)
	r.Use(middleware.ErrorHandler)
	r.Use(middleware.CORS)
	r.Use(gin.Recovery())
	r.Use(middleware.RecoverAbortedRequest())

	// register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("testcustom", customvalidator.TestCustom)
	}

	// general prefix for all routes
	v1 := r.Group("api/v1")

	// subgroups for nested grouping
	addTesting(v1)
	addContacts(v1)

	return r
}
