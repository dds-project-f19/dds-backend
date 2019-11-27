package controllers

import (
	"dds-backend/common"
	"dds-backend/database"
	"dds-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ManagerController struct {
	ControllerBase
}

// POST /manager/register_worker
// {"username":"required", "password":"required", "name":"", "surname":"", "phone":"", "address":""}
// 201: {"token":"1234567"}
// 400,409,500: {"message":"123"}
func (a *ManagerController) RegisterWorker(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Manager))
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
		newUser.Claim = common.Worker
		newUser.GameType = auth.GameType
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

// GET /manager/list_workers
// HEADERS: {Authorization: token}
// {}
// 200: {"users":[{"username":""...}]}
// 401,500: {"message":"123"}
func (a *ManagerController) ListWorkers(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Manager))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	searchItem := models.User{GameType: auth.GameType}
	var users []models.User
	resp := database.DB.Model(&models.User{}).Where(&searchItem).Find(&users)
	if err := resp.Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, resp.Error.Error())
		return
	}
	var toDump []interface{}
	for _, elem := range users {
		toDump = append(toDump, elem.ToMap())
	}
	a.JsonSuccess(c, http.StatusOK, gin.H{"users": toDump})
}

// DEPRECATED
// /manager/get_worker/{username}
//func (a *WorkerController) GetWorker(c *gin.Context) {
//	_, err := CheckAuthConditional(c, HasEqualOrHigherClaim(Manager))
//	if err != nil {
//		a.JsonFail(c, http.StatusUnauthorized, err.Error())
//		return
//	}
//
//	user := models.User{Username: c.Param("username")}
//	if database.DB.Model(&models.User{}).Where(&user).First(&user).RecordNotFound() {
//		a.JsonFail(c, http.StatusNotFound, "user not found")
//		return
//	}
//	a.JsonSuccess(c, http.StatusOK, user.ToMap())
//}

// DELETE /manager/remove_worker/{username}
// HEADERS: {Authorization: token}
// {}
// 200: {}
// 400,401,404,500: {"message":"123"}
func (a *ManagerController) RemoveWorker(c *gin.Context) {
	_, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Manager))
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
		// TODO: add logs and don't send error messages to users
		a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
		return
	}

	if err := database.DB.Model(&models.User{}).Delete(&user).Error; err != nil {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
		return
	}
	a.JsonSuccess(c, http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// PATCH /manager/add_available_items
// HEADERS: {Authorization: token}
// {"itemtype":"123","count":77}
// 200: {}
// 400,401,500: {"message":"123"}
func (a *ManagerController) AddAvailableItems(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Manager))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	var item models.AvailableItem
	if err := c.Bind(&item); err == nil {
		tx := database.DB.Begin()

		searchItem := models.AvailableItem{GameType: auth.GameType}
		item.GameType = auth.GameType

		res := tx.Model(&models.AvailableItem{}).Where(&searchItem).First(&searchItem)
		if res.RecordNotFound() {
			res := tx.Model(&models.AvailableItem{}).Create(&item)
			if res.Error != nil {
				tx.Rollback()
				a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
				return
			}
		} else if res.Error != nil {
			tx.Rollback()
			a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
			return
		} else {
			searchItem.Count += item.Count
			res = tx.Model(&models.AvailableItem{}).Save(&searchItem)
			if res.Error != nil {
				tx.Rollback()
				a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
				return
			}
		}
		if err := tx.Commit().Error; err != nil {
			a.JsonFail(c, http.StatusInternalServerError, err.Error())
			return
		}
		a.JsonSuccess(c, http.StatusOK, gin.H{})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}

}

// PATCH /manager/remove_available_items
// HEADERS: {Authorization: token}
// {"itemtype":"123","count":77}
// 200: {}
// 400,401,500: {"message":"123"}
func (a *ManagerController) RemoveAvailableItems(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Manager))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	var item models.AvailableItem
	if err := c.Bind(&item); err == nil {
		tx := database.DB.Begin()

		searchItem := models.AvailableItem{ItemType: item.ItemType, GameType: auth.GameType}
		item.GameType = auth.GameType
		res := tx.Model(&models.AvailableItem{}).Where(&searchItem).First(&searchItem)
		if res.RecordNotFound() {
			res := tx.Model(&models.AvailableItem{}).Create(&item)
			if res.Error != nil {
				tx.Rollback()
				a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
				return
			}
		} else if res.Error != nil {
			tx.Rollback()
			a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
			return
		} else {
			if searchItem.Count < item.Count {
				a.JsonFail(c, http.StatusBadRequest, "not enough items for deletion")
				return
			}
			searchItem.Count -= item.Count
			res = tx.Model(&models.AvailableItem{}).Save(&searchItem)
			if res.Error != nil {
				tx.Rollback()
				a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
				return
			}
		}
		if err := tx.Commit().Error; err != nil {
			a.JsonFail(c, http.StatusInternalServerError, err.Error())
			return
		}
		a.JsonSuccess(c, http.StatusOK, gin.H{})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

// GET /manager/list_available_items
// HEADERS: {Authorization: token}
// {}
// 200: {"items":[{"itemtype":"123","count":77}]}
// 401,500: {"message":"123"}
func (a *ManagerController) ListAvailableItems(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Manager))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	searchItem := models.AvailableItem{GameType: auth.GameType}
	var availableItems []models.AvailableItem

	resp := database.DB.Model(&models.AvailableItem{}).Where(&searchItem).Find(&availableItems)
	if err := resp.Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, resp.Error.Error())
		return
	}
	var toDump []interface{}
	for _, elem := range availableItems {
		if elem.Count > 0 {
			toDump = append(toDump, elem.ToMap())
		}
	}
	a.JsonSuccess(c, http.StatusOK, gin.H{"items": toDump})
	// list available items in form {"itemtype1":{"count":123}, ...}
}

// GET /manager/list_taken_items
// HEADERS: {Authorization: token}
// {}
// 200: {"items":[{"takenby":"username","itemtype":"123","assignedtoslot":"123"}]}
// 401,500: {"message":"123"}
func (a *ManagerController) ListTakenItems(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Manager))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	searchItem := models.AvailableItem{GameType: auth.GameType}
	var takenItems []models.TakenItem
	resp := database.DB.Model(&models.TakenItem{}).Where(&searchItem).Find(&takenItems)
	if err := resp.Error; err != nil {
		a.JsonFail(c, http.StatusInternalServerError, resp.Error.Error())
		return
	}
	var toDump []interface{}
	for _, elem := range takenItems {
		toDump = append(toDump, elem.ToMap())
	}
	a.JsonSuccess(c, http.StatusOK, gin.H{"items": toDump})
	// list taken items in form {"itemtype1":{"takenby":"username1", }}
}

func (a *ManagerController) SetWorkerSchedule(c *gin.Context) {
	// TODO implement
}
