package controllers

import (
	"dds-backend/database"
	"dds-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WorkerController struct {
	ControllerBase
}

// POST /worker/login
// {"username":"123", "password":"456"}
// 200: {"token":"1234567"}
// 400,403: {"message":"123"}
func (a *WorkerController) Login(c *gin.Context) {
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

// POST /worker/register
// {"username":"required", "password":"required", "name":"", "surname":"", "phone":"", "address":""}
// 201: {"token":"1234567"}
// 400,409,500: {"message":"123"}
func (a *WorkerController) Register(c *gin.Context) {
	var newUser models.User
	if err := c.Bind(&newUser); err == nil {
		if valid, msg := newUser.IsValid(); !valid {
			a.JsonFail(c, http.StatusBadRequest, msg)
			return
		}
		newUser.Claim = Worker
		tx := database.DB.Begin()
		existingUser := models.User{Username: newUser.Username}
		res := tx.Model(&models.User{}).Where(&existingUser).First(&existingUser)
		if res.RecordNotFound() {
			newUser.Password = Hash(newUser.Password)
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

// GET /worker/get
// HEADERS: {Authorization: token}
// {}
// 200: {"username":"required", "name":"", "surname":"", "phone":"", "address":""}
// 401,404: {"message":"123"}
func (a *WorkerController) Get(c *gin.Context) {
	auth, err := CheckAuthConditional(c)
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
	auth, err := CheckAuthConditional(c)
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
	if _, err := CheckAuthConditional(c); err != nil {
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
	auth, err := CheckAuthConditional(c, HasEqualOrHigherClaim(Worker))
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

		tx := database.DB.Begin()
		res := tx.Model(&models.AvailableItem{}).Where(&available).First(&available)
		if res.Error != nil {
			a.JsonFail(c, http.StatusBadRequest, res.Error.Error())
			return
		}
		taken.ItemType = available.ItemType
		taken.TakenBy = auth.Username
		taken.AssignedToSlot = request.Slot

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
	auth, err := CheckAuthConditional(c, HasEqualOrHigherClaim(Worker))
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

		tx := database.DB.Begin()
		res := tx.Model(&models.TakenItem{}).Where(&taken).First(&taken)
		if res.Error != nil {
			a.JsonFail(c, http.StatusBadRequest, res.Error.Error())
			return
		}
		available.ItemType = request.ItemType
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
	_, err := CheckAuthConditional(c, HasEqualOrHigherClaim(Worker))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	var items []models.AvailableItem
	resp := database.DB.Model(&models.AvailableItem{}).Find(&items)
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
	auth, err := CheckAuthConditional(c, HasEqualOrHigherClaim(Worker))
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}

	searchItem := models.TakenItem{TakenBy: auth.Username}
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
