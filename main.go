package main

import (
	"dds-backend/config"
	"dds-backend/database"
	"dds-backend/routes"
	"fmt"
	"github.com/gin-contrib/cors"
	"time"
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

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		//AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Run(config.Get().Addr)
}
