package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HeartBeat(c *gin.Context) {
	res := JSONResponse{
		Response: nil,
		Message:  "heartbeat is alive",
	}
	c.JSON(http.StatusOK, res)
	go SendGmail(
		"daronphang@gmail.com",
		"daronphang@gmail.com",
		"test gmail api",
		"hello world!",
	)
}
