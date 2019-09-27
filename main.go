package main

import (
	"dds-backend/config"
	"dds-backend/database"
	"dds-backend/routes"
	"fmt"
)

func main() {
	if err := config.Load("config.yaml"); err != nil {
		fmt.Println("Failed to load configuration: " + err.Error())
		return
	}

	db, err := database.InitDB()
	if err != nil {
		fmt.Println("error opening database: " + err.Error())
		return
	}
	defer db.Close()

	router := routes.InitRouter()
	router.Run(config.Get().Addr)
}
