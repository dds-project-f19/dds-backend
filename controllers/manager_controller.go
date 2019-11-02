package controllers

//type Manager struct {
//	ControllerBase
//}
//
//func (a *Manager) ListUsers(c *gin.Context) {
//	var users []models.User
//	var request LoginRequest
//	var dbUser models.User
//
//	if err := c.ShouldBind(&request); err == nil {
//		password := []byte(request.Password)
//		ctx := sha256.New()
//		ctx.Write(password)
//		cipherStr := ctx.Sum(nil)
//		hexpass := hex.EncodeToString(cipherStr)
//
//		if database.DB.Where("username = ?", request.Username).First(&dbUser).Error != nil {
//			a.JsonFail(c, http.StatusNotFound, "User not found")
//			return
//		}
//
//		if hexpass != dbUser.Password {
//			a.JsonFail(c, http.StatusForbidden, "Wrong password")
//			return
//		}
//
//		if dbUser.Username != "admin" {
//			a.JsonFail(c, http.StatusForbidden, "Unauthorized")
//			return
//		}
//
//		database.DB.Select("*").Order("id").Find(&users)
//		a.JsonSuccess(c, http.StatusOK, gin.H{"data": users})
//	} else {
//		a.JsonFail(c, http.StatusBadRequest, err.Error())
//	}
//
//}
//
//func (a *Manager) Login(c *gin.Context) {
//	var request LoginRequest
//	var dbUser models.User
//
//	if err := c.ShouldBind(&request); err == nil {
//		password := []byte(request.Password)
//		ctx := sha256.New()
//		ctx.Write(password)
//		cipherStr := ctx.Sum(nil)
//		hexpass := hex.EncodeToString(cipherStr)
//
//		if database.DB.Where("username = ?", request.Username).First(&dbUser).Error != nil {
//			a.JsonFail(c, http.StatusNotFound, "User not found")
//			return
//		}
//
//		if hexpass != dbUser.Password {
//			a.JsonFail(c, http.StatusForbidden, "Wrong password")
//			return
//		}
//
//		a.JsonSuccess(c, http.StatusOK, gin.H{})
//	} else {
//		a.JsonFail(c, http.StatusBadRequest, err.Error())
//	}
//}
//
//func (a *Manager) Store(c *gin.Context) {
//	var request CreateRequest
//
//	if err := c.ShouldBind(&request); err == nil {
//		var count int
//		database.DB.Model(&models.User{}).Where("username = ?", request.Username).Count(&count)
//
//		if count > 0 {
//			a.JsonFail(c, http.StatusBadRequest, "Username already exists")
//			return
//		}
//
//		password := []byte(request.Password)
//		ctx := sha256.New()
//		ctx.Write(password)
//		cipherStr := ctx.Sum(nil)
//		user := models.User{
//			Username: request.Username,
//			Name:     request.Name,
//			Surname:  request.Surname,
//			Phone:    request.Phone,
//			Address:  request.Address,
//			Password: hex.EncodeToString(cipherStr),
//		}
//
//		if err := database.DB.Create(&user).Error; err != nil {
//			a.JsonFail(c, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		a.JsonSuccess(c, http.StatusCreated, gin.H{"message": "User created successfully"})
//	} else {
//		a.JsonFail(c, http.StatusBadRequest, err.Error())
//	}
//}
//
//func (a *Manager) Update(c *gin.Context) {
//	var request UpdateRequest
//
//	if err := c.ShouldBind(&request); err == nil {
//		var user models.User
//		if database.DB.First(&user, c.Param("id")).Error != nil {
//			a.JsonFail(c, http.StatusNotFound, "User not found")
//			return
//		}
//
//		user.Name = request.Name
//
//		if err := database.DB.Save(&user).Error; err != nil {
//			a.JsonFail(c, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		a.JsonSuccess(c, http.StatusCreated, gin.H{})
//	} else {
//		a.JsonFail(c, http.StatusBadRequest, err.Error())
//	}
//}
//
//func (a *Manager) Show(c *gin.Context) {
//	var user models.User
//
//	if database.DB.Select("id, name, username, created_at, updated_at").First(&user, c.Param("id")).Error != nil {
//		a.JsonFail(c, http.StatusNotFound, "User not found")
//		return
//	}
//
//	a.JsonSuccess(c, http.StatusCreated, gin.H{"data": user})
//}
//
//func (a *Manager) Destroy(c *gin.Context) {
//	var user models.User
//
//	if database.DB.First(&user, c.Param("id")).Error != nil {
//		a.JsonFail(c, http.StatusNotFound, "User not found")
//		return
//	}
//
//	if err := database.DB.Unscoped().Delete(&user).Error; err != nil {
//		a.JsonFail(c, http.StatusBadRequest, err.Error())
//		return
//	}
//
//	a.JsonSuccess(c, http.StatusCreated, gin.H{})
//
//}
//
//type UpdateRequest struct {
//	Name string `form:"name" json:"name" binding:"required"`
//}
//
//type CreateRequest struct {
//	Username string `form:"username" json:"username" binding:"required"`
//	Name     string `form:"name" json:"name" binding:"required"`
//	Surname  string `form:"surname" json:"surname" binding:"required"`
//	Phone    string `form:"phone" json:"phone" binding:"required"`
//	Address  string `form:"address" json:"address" binding:"required"`
//	Password string `form:"password" json:"password" binding:"required"`
//}
//
//type LoginRequest struct {
//	Username string `form:"username" json:"username" binding:"required"`
//	Password string `form:"password" json:"password" binding:"required"`
//}
