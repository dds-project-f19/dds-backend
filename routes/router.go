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
		worker := new(controllers.WorkerController)
		workers.POST("/login", worker.Login)
		workers.POST("/register", worker.Register)
		workers.GET("/get", worker.Get)
		workers.PATCH("/update", worker.Update)
		workers.GET("/check_access", worker.CheckAccess)
		workers.POST("/take_item", worker.TakeItem)
		workers.POST("/return_item", worker.ReturnItem)
		workers.GET("/available_items", worker.AvailableItems)
		workers.GET("/taken_items", worker.TakenItems)
	}

	managers := router.Group("/manager")
	{
		manager := new(controllers.ManagerController)
		managers.POST("/login", manager.Login)
		managers.GET("/list_workers", manager.ListWorkers)
		managers.DELETE("/remove_worker/:username", manager.RemoveWorker)
		managers.PATCH("/add_available_items", manager.AddAvailableItems)
		managers.PATCH("/remove_available_items", manager.RemoveAvailableItems)
		managers.GET("/list_available_items", manager.ListAvailableItems)
		managers.GET("/list_taken_items", manager.ListTakenItems)

	}

	admins := router.Group("/admin")
	{
		admin := new(controllers.AdminController)
		admins.POST("/register_manager", admin.RegisterManager)
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
