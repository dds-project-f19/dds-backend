package controllers

import (
	"dds-backend/common"
	"dds-backend/database"
	"dds-backend/models"
	"dds-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
		newUser.Claim = common.Worker
		newUser.GameType = auth.GameType
		if valid, msg := newUser.IsValid(); !valid {
			a.JsonFail(c, http.StatusBadRequest, msg)
			return
		}
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
		a.JsonSuccess(c, http.StatusCreated, gin.H{})
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

	searchItem := models.User{GameType: auth.GameType, Claim: common.Worker}
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

// PATCH /manager/set_available_items
// HEADERS: {Authorization: token}
// {"itemtype":"123","count":77}
// 200: {}
// 400,401,500: {"message":"123"}
func (a *ManagerController) SetAvailableItems(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Manager))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	var item models.AvailableItem
	if err := c.Bind(&item); err == nil {
		if err := item.CheckValid(); err != nil { // validate
			a.JsonFail(c, http.StatusBadRequest, err.Error())
		}

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
			searchItem.Count = item.Count
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

	searchItem := models.TakenItem{GameType: auth.GameType}
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

// POST /manager/set_worker_schedule
// HEADERS: {Authorization: token}
// {"username":"abc", "starttime":"12:13", "endtime":"14:13", "workdays";"1,4,5"} - workdays are monday, thursday, friday
// 200: {}
// 400, 401, 404, 500: {"message":"123"}
func (a *ManagerController) SetWorkerSchedule(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Manager))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}
	type ScheduleRequest struct {
		Username  string `binding:"required"`
		StartTime string `binding:"required"`
		EndTime   string `binding:"required"`
		Workdays  string `binding:"required"`
	}
	request := ScheduleRequest{}

	if err := c.Bind(&request); err == nil {
		t1, err := models.LoadTimePoint(request.StartTime)
		if err != nil || !t1.IsValid() {
			a.JsonFail(c, http.StatusBadRequest, "start time ill-formed")
			return
		}
		t2, err := models.LoadTimePoint(request.EndTime)
		if err != nil || !t2.IsValid() {
			a.JsonFail(c, http.StatusBadRequest, "end time ill-formed")
			return
		}
		wks, err := models.LoadWeekdays(request.Workdays)
		if err != nil {
			a.JsonFail(c, http.StatusBadRequest, "workdays ill-formed")
			return
		}
		if t2.Before(t1) {
			a.JsonFail(c, http.StatusBadRequest, "time ordering mismatch")
			return
		}
		err = services.SetSchedule(request.Username, auth.GameType, wks, t1, t2)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				a.JsonFail(c, http.StatusNotFound, err.Error())
				return
			} else {
				a.JsonFail(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		a.JsonSuccess(c, http.StatusOK, gin.H{})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

// GET /manager/get_worker_schedule/{username}
// HEADERS: {Authorization: token}
// {}
// 200: {"starttime":"12:13", "endtime":"14:13", "workdays";"1,4,5"}
// 401, 404, 500: {"message":"123"}
func (a *ManagerController) GetWorkerSchedule(c *gin.Context) {
	_, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Manager))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}
	t1, t2, wks, err := services.GetSchedule(c.Param("username"))
	if err != nil {
		if _, ok := err.(*services.ScheduleNotFoundError); ok {
			a.JsonFail(c, http.StatusNotFound, err.Error())
		} else {
			a.JsonFail(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		a.JsonSuccess(c, http.StatusOK, gin.H{"starttime": t1.ToStr(), "endtime": t2.ToStr(),
			"workdays": models.StoreWeekdays(wks)})
	}
}

// POST /manager/check_overlap
// HEADERS: {Authorization: token}
// {"username":"worker1", "starttime":"10:20", "endtime":"10:30", "workdays":"1,2,3"}
// 200: {"overlap":true} - true (not string) for overlap error, false for no overlap
// 401, 404, 500: {"message":"123"}
func (a *ManagerController) CheckOverlap(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c, common.HasEqualOrHigherClaim(common.Manager))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}
	type CheckRequest struct {
		Username  string
		StartTime string `binding:"required"`
		EndTime   string `binding:"required"`
		Workdays  string
	}
	var request CheckRequest
	var schs []models.UserSchedule
	searchItem := models.UserSchedule{GameType: auth.GameType}
	if err := c.Bind(&request); err == nil {
		res := database.DB.Model(&models.UserSchedule{}).Where(&searchItem).Find(&schs)
		if res.Error != nil {
			a.JsonFail(c, http.StatusInternalServerError, res.Error.Error())
			return
		}
		t1cur, err := models.LoadTimePoint(request.StartTime)
		if err != nil {
			a.JsonFail(c, http.StatusBadRequest, "parsing time failed")
			return
		}
		t2cur, err := models.LoadTimePoint(request.EndTime)
		if err != nil {
			a.JsonFail(c, http.StatusBadRequest, "parsing time failed")
			return
		}
		wrk, err := models.LoadWeekdays(request.Workdays)
		if err != nil {
			a.JsonFail(c, http.StatusBadRequest, "parsing weekdays failed")
			return
		}
		for _, e := range schs {
			if e.Username == request.Username {
				continue // don't check for same user
			}
			t1, err := models.LoadTimePoint(e.StartTime)
			if err != nil {
				a.JsonFail(c, http.StatusInternalServerError, "parsing time failed")
				return
			}
			t2, err := models.LoadTimePoint(e.EndTime)
			if err != nil {
				a.JsonFail(c, http.StatusInternalServerError, "parsing time failed")
				return
			}
			we, err := models.LoadWeekdays(e.Workdays)
			if err != nil {
				a.JsonFail(c, http.StatusInternalServerError, "parsing workdays failed")
				return
			}
			if t1cur.Before(t2) && t1.Before(t1cur) || t2cur.Before(t2) && t1.Before(t2cur) ||
				t1.Before(t2cur) && t1cur.Before(t1) || t2.Before(t2cur) && t1cur.Before(t2) {
				for _, w1 := range we {
					for _, w2 := range wrk {
						if w1 == w2 {
							a.JsonSuccess(c, http.StatusOK, gin.H{"overlap": true})
							return
						}
					}
				}
			}

		}
		a.JsonSuccess(c, http.StatusOK, gin.H{"overlap": false})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}
