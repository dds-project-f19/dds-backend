package controllers

import (
	"crypto/sha256"
	"dds-backend/database"
	"dds-backend/models"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	ControllerBase
}

func (a *User) ListUsers(c *gin.Context) {
	var users []models.User
	var dbUser models.User

	if err := c.ShouldBind(&dbUser); err == nil {
		password := []byte(dbUser.Password)
		ctx := sha256.New()
		ctx.Write(password)
		cipherStr := ctx.Sum(nil)
		hexpass := hex.EncodeToString(cipherStr)

		if database.DB.Where("username = ?", dbUser.Username).First(&dbUser).Error != nil {
			a.JsonFail(c, http.StatusNotFound, "User not found")
			return
		}

		if hexpass != dbUser.Password {
			a.JsonFail(c, http.StatusForbidden, "Wrong password")
			return
		}

		if dbUser.Username != "admin" {
			a.JsonFail(c, http.StatusForbidden, "Unauthorized")
			return
		}

		database.DB.Select("*").Order("id").Find(&users)
		a.JsonSuccess(c, http.StatusOK, gin.H{"data": users})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}

}

func (a *User) Login(c *gin.Context) {
	type RequestBody struct {
		Username string `binding:"required"`
		Password string `binding:"required"`
	}
	var body RequestBody

	if err := c.ShouldBind(&body); err == nil {

		//token, err := Authorize(body.Username, hexpass)

		a.JsonSuccess(c, http.StatusOK, gin.H{})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

func (a *User) Register(c *gin.Context) {
	var newUser models.User

	if err := c.Bind(&newUser); err == nil {
		tx := database.DB.Begin()
		existingUser := models.User{Username: newUser.Username}
		tx.Find(&existingUser)
		if tx.RecordNotFound() {

		} else {

		}

		//password := []byte(request.Password)
		//ctx := sha256.New()
		//ctx.Write(password)
		//cipherStr := ctx.Sum(nil)
		//user := models.User(request)

		//if err := database.DB.Create(&user).Error; err != nil {
		//	a.JsonFail(c, http.StatusBadRequest, err.Error())
		//	return
		//}

		a.JsonSuccess(c, http.StatusCreated, gin.H{"message": "User created successfully"})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

func (a *User) Update(c *gin.Context) {
	var request models.User

	if err := c.ShouldBind(&request); err == nil {
		var user models.User
		if database.DB.First(&user, c.Param("id")).Error != nil {
			a.JsonFail(c, http.StatusNotFound, "User not found")
			return
		}

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
