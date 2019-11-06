package controllers

import (
	"dds-backend/database"
	"dds-backend/models"
	"fmt"
	"github.com/jinzhu/gorm"
)

func InitializeDefaultUsers() {
	admin := models.User{
		Model:    gorm.Model{},
		Username: "admin",
		Password: "password",
		Name:     "Maksim",
		Surname:  "Surkov",
		Phone:    "123",
		Address:  "Github str. 1, nil",
		Claim:    Admin,
	}
	if err := database.DB.Create(&admin).Error; err != nil {
		fmt.Print("admin already exists")
	}
}
