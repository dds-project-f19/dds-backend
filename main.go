package main

import (
	"dds-backend/config"
	"dds-backend/database"
	"dds-backend/routes"
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

	router := routes.InitRouter()

	router.Run(generalConfig.Address)
}
