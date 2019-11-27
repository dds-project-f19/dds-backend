package controllers

import (
	"dds-backend/common"
	"dds-backend/database"
	"dds-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WorkerController struct {
	ControllerBase
}

// GET /worker/get
// HEADERS: {Authorization: token}
// {}
// 200: {"username":"required", "name":"", "surname":"", "phone":"", "address":""}
// 401,404: {"message":"123"}
func (a *WorkerController) Get(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c)
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	user := models.User{Username: auth.Username}
	if database.DB.Model(&models.User{}).Where(&user).First(&user).RecordNotFound() {
		a.JsonFail(c, http.StatusNotFound, "user not found")
		return
	}
	a.JsonSuccess(c, http.StatusOK, user.ToMap())
}

// PATCH /worker/update
// HEADERS: {Authorization: token}
// {"username":"required", "name":"", "surname":"", "phone":"", "address":""}
// 200: {}
// 400,401,404: {"message":"123"}
// TODO: fix gorm requests and decide on update semantics
func (a *WorkerController) Update(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c)
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}
	var request models.User

	if err := c.ShouldBind(&request); err == nil {

		user := models.User{Username: auth.Username}
		if err := database.DB.Find(&user).Error; err != nil {
			a.JsonFail(c, http.StatusNotFound, err.Error())
			return
		}
		user = request

		if err := database.DB.Save(&user).Error; err != nil {
			a.JsonFail(c, http.StatusBadRequest, err.Error())
			return
		}

		a.JsonSuccess(c, http.StatusOK, gin.H{})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

func (a *WorkerController) CheckAccess(c *gin.Context) {
	if _, err := common.CheckAuthConditional(c); err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
	} else {
		a.JsonSuccess(c, http.StatusOK, gin.H{"message": "you have access to perform this call"})
	}
}

// POST /worker/take_item
// HEADERS: {Authorization: token}
// {"itemtype":"123", "slot":"123"}
// 201: {"message":"request done, blah blah"}
// 400,401,500: {"message":"123"}
func (a *WorkerController) TakeItem(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Worker))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}
	type moveRequest struct {
		ItemType string
		Slot     string
	}

	var request moveRequest
	var available models.AvailableItem
	var taken models.TakenItem

	if err := c.Bind(&request); err == nil {
		available.ItemType = request.ItemType
		available.GameType = auth.GameType

		tx := database.DB.Begin()
		res := tx.Model(&models.AvailableItem{}).Where(&available).First(&available)
		if res.Error != nil {
			a.JsonFail(c, http.StatusBadRequest, res.Error.Error())
			return
		}
		taken.ItemType = available.ItemType
		taken.TakenBy = auth.Username
		taken.AssignedToSlot = request.Slot
		taken.GameType = auth.GameType

		if available.Count <= 0 {
			a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
			return
		}
		available.Count--
		res = tx.Model(&models.AvailableItem{}).Save(&available)
		if res.Error != nil {
			tx.Rollback()
			a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		res = tx.Model(&models.TakenItem{}).Create(&taken)
		if res.Error != nil {
			tx.Rollback()
			a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
			return
		}
		res = tx.Commit()
		if res.Error != nil {
			a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
			return
		}
		a.JsonSuccess(c, http.StatusCreated, gin.H{"message": "item moved successfully"})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

// POST /worker/return_item
// HEADERS: {Authorization: token}
// {"itemtype":"123", "slot":"123"}
// 201: {"message":"request done, blah blah"}
// 400,401,500: {"message":"123"}
func (a *WorkerController) ReturnItem(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Worker))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}
	type moveRequest struct {
		ItemType string
		Slot     string
	}

	var request moveRequest
	var available models.AvailableItem
	var taken models.TakenItem

	if err := c.Bind(&request); err == nil {
		taken.ItemType = request.ItemType
		taken.AssignedToSlot = request.Slot
		taken.TakenBy = auth.Username
		taken.GameType = auth.GameType

		tx := database.DB.Begin()
		res := tx.Model(&models.TakenItem{}).Where(&taken).First(&taken)
		if res.Error != nil {
			a.JsonFail(c, http.StatusBadRequest, res.Error.Error())
			return
		}
		available.ItemType = request.ItemType
		available.GameType = auth.GameType
		res = tx.Model(&models.AvailableItem{}).Where(&available).First(&available)
		if res.Error != nil {
			a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
			return
		}
		available.Count++
		res = tx.Model(&models.AvailableItem{}).Save(&available)
		if res.Error != nil {
			tx.Rollback()
			a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		res = tx.Model(&models.TakenItem{}).Delete(taken)
		if res.Error != nil {
			tx.Rollback()
			a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
			return
		}
		res = tx.Commit()
		if res.Error != nil {
			a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
			return
		}
		a.JsonSuccess(c, http.StatusCreated, gin.H{"message": "item moved successfully"})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

// GET /worker/list_available_items
// HEADERS: {Authorization: token}
// {}
// 200: {"items":[{"itemtype":"123","count":77}]}
// 401,500: {"message":"123"}
func (a *WorkerController) AvailableItems(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Worker))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	searchItem := models.AvailableItem{GameType: auth.GameType}
	var items []models.AvailableItem
	resp := database.DB.Model(&models.AvailableItem{}).Where(&searchItem).Find(&items)
	if err := resp.Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, resp.Error.Error())
		return
	}
	var toDump []interface{}
	for _, elem := range items {
		if elem.Count > 0 {
			toDump = append(toDump, elem.ToMap())
		}
	}
	a.JsonSuccess(c, http.StatusOK, gin.H{"items": toDump})
}

// GET /worker/list_taken_items
// HEADERS: {Authorization: token}
// {}
// 200: {"items":[{"takenby":"username","itemtype":"123","assignedtoslot":"123"}]}
// 401,500: {"message":"123"}
func (a *WorkerController) TakenItems(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Worker))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	searchItem := models.TakenItem{TakenBy: auth.Username, GameType: auth.GameType}
	var items []models.TakenItem
	resp := database.DB.Model(&models.TakenItem{}).Where(&searchItem).Find(&items)
	if err := resp.Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, resp.Error.Error())
		return
	}
	var toDump []interface{}
	for _, elem := range items {
		toDump = append(toDump, elem.ToMap())
	}
	a.JsonSuccess(c, http.StatusOK, gin.H{"items": toDump})
}
