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

func (a *ManagerController) Login(c *gin.Context) {
	// same as user login
}

// TODO: fix gorm requests
func (a *ManagerController) ListWorkers(c *gin.Context) {
	if _, err := CheckAuthConditional(c, HasEqualOrHigherClaim(Manager)); err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	// TODO: test for no users
	var users []models.User
	resp := database.DB.Model(&models.User{}).Find(&users)
	if err := resp.Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, resp.Error.Error())
	}
}

func (a *WorkerController) GetWorker(c *gin.Context) {

}

// TODO: fix gorm requests
func (a *WorkerController) RemoveWorker(c *gin.Context) {
	_, err := CheckAuthConditional(c, HasEqualOrHigherClaim(Manager))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	user := models.User{Username: c.Param("username")}

	res := database.DB.Model(&models.User{}).Where(&user).First(&user)
	if res.RecordNotFound() {
		a.JsonFail(c, http.StatusNotFound, "worker not found")
		return
	} else if res.Error != nil {
		// TODO: add logs and close error interface to users
		a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
	}

	if err := database.DB.Model(&models.User{}).Delete(&user).Error; err != nil {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (a *WorkerController) AddAvailableItems(c *gin.Context) {

}

func (a *WorkerController) RemoveAvailableItems(c *gin.Context) {

}

func (a *WorkerController) GetAvailableItems(c *gin.Context) {

}

func (a *WorkerController) ListTakenItems(c *gin.Context) {

}
