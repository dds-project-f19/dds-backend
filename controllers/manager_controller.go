package controllers

import (
	"dds-backend/database"
	"dds-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ManagerController struct {
	ControllerBase
}

// TODO: fix gorm requests
func (a *ManagerController) ListUsers(c *gin.Context) {
	if err := checkAuthConditional(c, HasEqualOrHigherClaim(Manager)); err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	var users []models.User
	resp := database.DB.Model(&models.User{}).Find(&users)
	if err := resp.Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, resp.Error.Error())
	}
}

// TODO: fix gorm requests
func (a *WorkerController) Destroy(c *gin.Context) {
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
