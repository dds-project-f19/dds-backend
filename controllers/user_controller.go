package controllers

import (
	"crypto/md5"
	"dds-backend/database"
	"dds-backend/models"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	ControllerBase
}

func (a *User) Index(c *gin.Context) {
	var users []models.User

	database.DB.Select("id, name, username, created_at, updated_at").Order("id").Find(&users)

	a.JsonSuccess(c, http.StatusOK, gin.H{"data": users})
}

func (a *User) Store(c *gin.Context) {
	var request CreateRequest

	if err := c.ShouldBind(&request); err == nil {
		var count int
		database.DB.Model(&models.User{}).Where("username = ?", request.Username).Count(&count)

		if count > 0 {
			a.JsonFail(c, http.StatusBadRequest, "Username already exists")
			return
		}

		password := []byte(request.Password)
		md5Ctx := md5.New()
		md5Ctx.Write(password)
		cipherStr := md5Ctx.Sum(nil)
		user := models.User{
			Username: request.Username,
			Name:     request.Name,
			Password: hex.EncodeToString(cipherStr),
		}

		if err := database.DB.Create(&user).Error; err != nil {
			a.JsonFail(c, http.StatusBadRequest, err.Error())
			return
		}

		a.JsonSuccess(c, http.StatusCreated, gin.H{"message": "User created successfully"})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

func (a *User) Update(c *gin.Context) {
	var request UpdateRequest

	if err := c.ShouldBind(&request); err == nil {
		var user models.User
		if database.DB.First(&user, c.Param("id")).Error != nil {
			a.JsonFail(c, http.StatusNotFound, "User not found")
			return
		}

		user.Name = request.Name

		if err := database.DB.Save(&user).Error; err != nil {
			a.JsonFail(c, http.StatusBadRequest, err.Error())
			return
		}

		a.JsonSuccess(c, http.StatusCreated, gin.H{})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

func (a *User) Show(c *gin.Context) {
	var user models.User

	if database.DB.Select("id, name, username, created_at, updated_at").First(&user, c.Param("id")).Error != nil {
		a.JsonFail(c, http.StatusNotFound, "User not found")
		return
	}

	a.JsonSuccess(c, http.StatusCreated, gin.H{"data": user})
}

func (a *User) Destroy(c *gin.Context) {
	var user models.User

	if database.DB.First(&user, c.Param("id")).Error != nil {
		a.JsonFail(c, http.StatusNotFound, "User not found")
		return
	}

	if err := database.DB.Unscoped().Delete(&user).Error; err != nil {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}

	a.JsonSuccess(c, http.StatusCreated, gin.H{})

}

type UpdateRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type CreateRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
