package main

import (
	"dds-backend/config"
	"dds-backend/controllers"
	"dds-backend/database"
	"dds-backend/routes"
	telegram_bot "dds-backend/services"
	"fmt"
)

func main() {
	currentConfig := config.LoadConfigFromCmdArgs()
	generalConfig := config.GetDefaultGeneralConfig()

	db, err := database.InitDB(currentConfig, generalConfig)
	if err != nil {
		fmt.Println("error opening database: " + err.Error())
		return
	}
	defer db.Close()

	controllers.InitializeDefaultUsers() // create user `admin`

	go telegram_bot.LaunchBot()

	router := routes.InitRouter()
	router.Run(generalConfig.Address) // run backend router
}
