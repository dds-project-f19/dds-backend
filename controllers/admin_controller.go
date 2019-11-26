package controllers

import (
	"dds-backend/common"
	"dds-backend/database"
	"dds-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminController struct {
	ControllerBase
}

// POST /admin/register_manager
// {"username":"required", "password":"required", "name":"", "surname":"", "phone":"", "address":""}
// 201: {}
// 400,401,409,500: {"message":"123"}
func (a *AdminController) RegisterManager(c *gin.Context) {
	_, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Admin))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	var newUser models.User
	if err := c.Bind(&newUser); err == nil {
		if valid, msg := newUser.IsValid(); !valid {
			a.JsonFail(c, http.StatusBadRequest, msg)
			return
		}
		newUser.Claim = common.Manager
		tx := database.DB.Begin()
		existingUser := models.User{Username: newUser.Username}
		res := tx.Model(&models.User{}).Where(&existingUser).First(&existingUser)
		if res.RecordNotFound() {
			newUser.Password = common.Hash(newUser.Password)
			err = tx.Create(&newUser).Error
			if err != nil {
				tx.Rollback()
				a.JsonFail(c, http.StatusBadRequest, err.Error())
				return
			}
		} else if res.Error != nil {
			tx.Rollback()
			a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
			return
		} else {
			tx.Rollback()
			a.JsonFail(c, http.StatusConflict, "user already exists")
			return
		}
		if err := tx.Commit().Error; err != nil {
			a.JsonFail(c, http.StatusInternalServerError, err.Error())
			return
		}
		a.JsonSuccess(c, http.StatusCreated, gin.H{"message": "user created successfully"})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}
