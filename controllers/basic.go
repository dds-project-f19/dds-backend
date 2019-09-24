package controllers

import (
	"github.com/gin-gonic/gin"
)

type ControllerBase struct {
}

func (basic *ControllerBase) JsonSuccess(c *gin.Context, status int, h gin.H) {
	h["status"] = "success"
	c.JSON(status, h)
	return
}

func (basic *ControllerBase) JsonFail(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status":  "fail",
		"message": message,
	})
}
