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
		Password: Hash("password"),
		Name:     "Maksim",
		Surname:  "Surkov",
		Phone:    "123",
		Address:  "Github str. 1, nil",
		Claim:    Admin,
	}
	existingAdmin := models.User{Username: admin.Username}
	if database.DB.Model(&models.User{}).Where(&existingAdmin).First(&existingAdmin).RecordNotFound() {
		if err := database.DB.Create(&admin).Error; err != nil {
			panic(err)
		}
	} else {
		fmt.Println("*** Admin already exists ***")
	}

}
