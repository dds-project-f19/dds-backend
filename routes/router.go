package routes

import (
	"dds-backend/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// setup CORS policy
	router.Use(cors.Default())

	users := router.Group("/users")
	{
		user := new(controllers.User)
		users.POST("/list", user.ListUsers)
		users.POST("/login", user.Login)
		users.POST("/register", user.Register)
		users.PATCH("/edit/:id", user.Update)
		users.DELETE("/remove/:id", user.Destroy)
		users.GET("/get/:id", user.Show)
	}
	game := router.Group("/inventory")
	{
		gameState := new(controllers.GameState)
		game.GET("/available", gameState.GetAvailableInventory) // available items for gametype x
		game.POST("/transfer", gameState.TransferInventory)     // from available inventory to slot of user y
		game.PATCH("/layout/edit", gameState.UpdateUserLayout)  // change layout of user y
		game.POST("/update", gameState.UpdateInventory)         // edit available items in inventory for gametype x
	}

	ping := new(controllers.Ping)
	router.GET("/ping", ping.Ping)

	// serve frontend source
	router.Use(static.Serve("/", static.LocalFile("./client/build", true)))

	return router

}
