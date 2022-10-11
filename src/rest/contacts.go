package rest

import (
	"portfolio_golang/src/controller"
	"portfolio_golang/src/middleware"

	"github.com/gin-gonic/gin"
)

func addContacts (rg *gin.RouterGroup) {
	contacts := rg.Group("/contacts")

	{
		contacts.GET(
			"",
			middleware.ConnectDB("MYSQL"),
			controller.ReadContacts,
		)
		contacts.GET(
			"/:email",
			middleware.ConnectDB("MYSQL"),
			controller.ReadContacts,
		)
		contacts.POST(
			"",
			middleware.ValidateSchema("CREATECONTACT"),
			middleware.ConnectDB("MYSQL"),
			controller.CreateContact,
		)
		contacts.PUT(
			"/:uid",
			middleware.ConnectDB("MYSQL"),
		)
		contacts.DELETE(
			"/:uid",
			middleware.ConnectDB("MYSQL"),
		)
	}
}

