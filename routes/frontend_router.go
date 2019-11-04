package routes

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func InitFrontendRouter() *gin.Engine {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./front/build", true)))
	return router
}
