package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Ping struct {
	ControllerBase
}

func (p *Ping) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "PONG")
}
