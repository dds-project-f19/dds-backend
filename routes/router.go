package routes

import (
	"dds-backend/config"
	"dds-backend/controllers"
	"dds-backend/database"
	"dds-backend/services"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// setup CORS policy
	router.Use(cors.Default())

	// TODO: setup logger

	apiGroup := router.Group("/api")

	commons := apiGroup.Group("/common")
	{
		common := new(controllers.CommonController)
		commons.POST("/login", common.Login)
		commons.GET("/telegram_join_link", common.GenerateTelegramJoinLink)

	}

	workers := apiGroup.Group("/worker")
	{
		worker := new(controllers.WorkerController)
		workers.GET("/get", worker.Get)
		workers.PATCH("/update", worker.Update)
		workers.GET("/check_access", worker.CheckAccess)
		workers.POST("/take_item", worker.TakeItem)
		workers.POST("/return_item", worker.ReturnItem)
		workers.GET("/list_available_items", worker.AvailableItems)
		workers.GET("/list_taken_items", worker.TakenItems)
	}

	managers := apiGroup.Group("/manager")
	{
		manager := new(controllers.ManagerController)
		managers.POST("/register_worker", manager.RegisterWorker)
		managers.GET("/list_workers", manager.ListWorkers)
		managers.DELETE("/remove_worker/:username", manager.RemoveWorker)
		managers.PATCH("/add_available_items", manager.AddAvailableItems)
		managers.PATCH("/remove_available_items", manager.RemoveAvailableItems)
		managers.GET("/list_available_items", manager.ListAvailableItems)
		managers.GET("/list_taken_items", manager.ListTakenItems)

	}

	admins := apiGroup.Group("/admin")
	{
		admin := new(controllers.AdminController)
		admins.POST("/register_manager", admin.RegisterManager)
	}

	// TODO: consider using decorators for access management
	// TODO: add claim checking when manager can delete another manager

	ping := new(controllers.Ping)
	{
		apiGroup.GET("/ping", ping.Ping)
	}

	router.Use(static.Serve("/", static.LocalFile("./front/build", true)))

	return router

}

// Must close database after use
func MakeServer() (*gin.Engine, config.GeneralConfig, *gorm.DB, error) {
	currentConfig := config.LoadConfigFromCmdArgs()
	generalConfig := config.GetDefaultGeneralConfig()

	db, err := database.InitDB(currentConfig, generalConfig)
	if err != nil {
		fmt.Println("error opening database: " + err.Error())
		return nil, generalConfig, db, err
	}
	controllers.InitializeDefaultUsers() // create user `admin`
	go services.LaunchBot()
	go services.LaunchScheduler()

	return InitRouter(), generalConfig, db, nil
}
