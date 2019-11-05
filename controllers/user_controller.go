package controllers

import (
	"dds-backend/database"
	"dds-backend/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type User struct {
	ControllerBase
}

func InitializeDefaultUsers() {
	admin := models.User{
		Model:    gorm.Model{},
		Username: "admin",
		Password: "password",
		Name:     "Maksim",
		Surname:  "Surkov",
		Phone:    "123",
		Address:  "Github str. 1, nil",
		Claim:    10,
	}
	if err := database.DB.Create(&admin).Error; err != nil {
		fmt.Print("admin already exists")
	}
}

// TODO: check and move to manager controller
func (a *User) ListUsers(c *gin.Context) {
	if err := checkAuthConditional(c, HasEqualOrHigherClaim(Manager)); err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	var users []models.User
	resp := database.DB.Find(&users)
	if err := resp.Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, resp.Error.Error())
	}
}

func (a *User) Login(c *gin.Context) {
	type RequestBody struct {
		Username string `binding:"required"`
		Password string `binding:"required"`
	}
	var request RequestBody

	if err := c.ShouldBind(&request); err == nil {
		token, err := Authorize(request.Username, Hash(request.Password))
		if err != nil {
			a.JsonFail(c, http.StatusForbidden, err.Error())
			return
		}
		a.JsonSuccess(c, http.StatusOK, gin.H{"token": token})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

func (a *User) Register(c *gin.Context) {
	var newUser models.User

	if err := c.Bind(&newUser); err == nil {
		tx := database.DB.Begin()
		existingUser := models.User{Username: newUser.Username}
		var count int
		tx.Model(&models.User{}).Where(&existingUser).Count(&count)
		if count <= 0 {
			newUser.Password = Hash(newUser.Password)
			err = tx.Create(&newUser).Error
			if err != nil {
				a.JsonFail(c, http.StatusBadRequest, err.Error())
				tx.Rollback()
				return
			}
		} else {
			a.JsonFail(c, http.StatusConflict, "user already exists")
			tx.Rollback()
			return
		}
		if err := tx.Commit().Error; err != nil {
			a.JsonFail(c, http.StatusInternalServerError, err.Error())
			return
		}
		token, err := Authorize(newUser.Username, newUser.Password) // note password is already hashed
		if err != nil {
			a.JsonFail(c, http.StatusInternalServerError, err.Error())
			return
		}
		a.JsonSuccess(c, http.StatusCreated, gin.H{"token": token, "message": "user created successfully"})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

// TODO: fix gorm requests
func (a *User) Update(c *gin.Context) {
	if err := checkAuthConditional(c, HasSameUsername(c.Param("username"))); err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}
	var request models.User

	if err := c.ShouldBind(&request); err == nil {

		user := models.User{Username: c.Param("username")}
		if err := database.DB.Find(&user).Error; err != nil {
			a.JsonFail(c, http.StatusNotFound, err.Error())
			return
		}
		user = request

		if err := database.DB.Save(&user).Error; err != nil {
			a.JsonFail(c, http.StatusBadRequest, err.Error())
			return
		}

		a.JsonSuccess(c, http.StatusCreated, gin.H{})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

// TODO: fix gorm requests
func (a *User) Get(c *gin.Context) {
	if err := checkAuthConditional(c, HasSameUsername(c.Param("username"))); err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	user := models.User{Username: c.Param("username")}
	if database.DB.Find(&user).RecordNotFound() {
		a.JsonFail(c, http.StatusNotFound, "user not found")
	}
	a.JsonSuccess(c, http.StatusCreated, user.ToMap())
}

// TODO: fix and move to manager controller
func (a *User) Destroy(c *gin.Context) {
	if err := checkAuthConditional(c, HasEqualOrHigherClaim(Manager)); err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	user := models.User{Username: c.Param("username")}
	if err := database.DB.Find(&user).Error; err != nil {
		a.JsonFail(c, http.StatusNotFound, err.Error())
		return
	}
	if err := database.DB.Delete(&user).Error; err != nil {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusOK, gin.H{"message": "user deleted successfully"})

}
