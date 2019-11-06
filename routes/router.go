package routes

import (
	"dds-backend/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// setup CORS policy
	router.Use(cors.Default())

	users := router.Group("/users")
	{
		user := new(controllers.WorkerController)
		users.POST("/login", user.Login)
		users.POST("/register", user.Register)
		users.PATCH("/edit/:username", user.Update)
		users.DELETE("/remove/:username", user.Destroy)
		users.GET("/get/:username", user.Get)
	}
	//game := router.Group("/inventory")
	//{
	//	//gameState := new(controllers.GameState)
	//	//game.GET("/available", gameState.GetAvailableInventory) // available items for gametype x
	//	//game.POST("/transfer", gameState.TransferInventory)     // from available inventory to slot of user y
	//	//game.PATCH("/layout/edit", gameState.UpdateUserLayout)  // change layout of user y
	//	//game.POST("/update", gameState.UpdateInventory)         // edit available items in inventory for gametype x
	//}

	ping := new(controllers.Ping)
	router.GET("/ping", ping.Ping)

	return router

}
