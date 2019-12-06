package controllers

import (
	"dds-backend/common"
	"dds-backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommonController struct {
	ControllerBase
}

// POST /common/login
// HEADERS: {}
// {"username":"123", "password":"456"}
// 200: {"token":"1234567", "claim":"worker|manager|admin", "gametype":"surgeon1"}
// 400,403: {"message":"123"}
func (a *CommonController) Login(c *gin.Context) {
	type RequestBody struct {
		Username string `binding:"required"`
		Password string `binding:"required"`
	}
	var request RequestBody

	if err := c.ShouldBind(&request); err == nil {
		auth, err := common.Authorize(request.Username, common.Hash(request.Password))
		if err != nil {
			a.JsonFail(c, http.StatusForbidden, err.Error())
			return
		}
		c.SetCookie("dds-auth-token", auth.Token, 60*60*12, "/", "", false, false)
		a.JsonSuccess(c, http.StatusOK, gin.H{"token": auth.Token, "claim": common.StringClaim(auth.Claim), "gametype": auth.GameType})
	} else {
		a.JsonFail(c, http.StatusBadRequest, err.Error())
	}
}

// GET /common/telegram_join_link
// HEADERS: {Authorization: token}
// {}
// 200: {"link":"t.me/bot_link/start=regkey123"}
// 401, 500: {"message":"123"}
func (a *CommonController) GenerateTelegramJoinLink(c *gin.Context) {
	auth, err := common.CheckAuthConditional(c)
	if err != nil {
		a.JsonFail(c, http.StatusUnauthorized, err.Error())
		return
	}
	link, err := services.GetChatRegistrationLink(auth.Username)
	if err != nil {
		a.JsonFail(c, http.StatusInternalServerError, err.Error())
	}
	a.JsonSuccess(c, http.StatusOK, gin.H{"link": link})
}
