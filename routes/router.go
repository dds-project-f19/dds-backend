package routes

import (
	"dds-backend/controllers"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	users := router.Group("/users")
	{
		user := new(controllers.User)
		users.GET("/list", user.Index)
		users.POST("/register", user.Store)
		users.PATCH("/edit/:id", user.Update)
		users.DELETE("/remove/:id", user.Destroy)
		users.GET("/get/:id", user.Show)

	}

	ping := new(controllers.Ping)
	router.GET("/ping", ping.Ping)

	return router

}
