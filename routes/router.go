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

	workers := router.Group("/worker")
	{
		user := new(controllers.WorkerController)
		workers.POST("/login", user.Login)
		workers.POST("/register", user.Register)
		workers.GET("/get", user.Get)
		workers.PATCH("/update", user.Update)
		workers.GET("/check_access", user.CheckAccess)
	}
	// TODO: consider using decorators for access management

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
