package controller

import (
	"fmt"
	"net/http"
	"net/mail"
	"sih/common"
	"sih/models"
	"sih/models/svrsn"
	"sih/pkg/utils"
	"strconv"
	"unicode/utf8"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthRequired is a simple middleware to check the session
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}

// siteOwnerRegister (POST)
func SiteOwnerRegister(c *gin.Context) {
	if utils.IsTrueVal(
		models.SvrSettingGet(svrsn.OwnerRegistered)) {
		c.String(http.StatusNotFound, "404 page not found")
	} else {
		var inputJSON SignupInfo
		if !registerReqDeal(c, &inputJSON) {
			return
		}

		id, _, errList := models.UserCreate(inputJSON.LoginName, inputJSON.ShowName, inputJSON.Email, inputJSON.Password, 1, false)

		if len(errList) == 0 {
			models.SvrSettingUpdate(svrsn.OwnerRegistered, strconv.FormatBool(true))
		}

		res := make(map[string]interface{})
		res["id"] = id
		res["errList"] = errList

		c.JSON(http.StatusOK, res)
		return
	}
}

func registerReqDeal(c *gin.Context, req *SignupInfo) bool {
	/*
		error message
		1: loginNameLen is used
		2: length of loginNameLen is not meet requirements
		3: length of showName is not meet requirements
		4: length of password is not meet requirements
		5: email is too long
		6: email not vaild
		7: email is used
		8: internal server error
	*/
	var errList []int
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request data"})
	}

	if models.IsLoginNameUsed(req.LoginName) {
		errList = append(errList, 1)
	}

	var loginNameLen = len(req.LoginName)
	if loginNameLen < 4 || loginNameLen > 30 {
		errList = append(errList, 2)
	}

	var showNameLen = utf8.RuneCountInString(req.ShowName)
	if showNameLen < 1 || showNameLen > 15 {
		errList = append(errList, 3)
	}

	if len(req.Password) < 6 {
		errList = append(errList, 4)
	}
	if len(req.Email) > 256 {
		errList = append(errList, 5)
	}

	_, err = mail.ParseAddress(req.Email)
	if len(req.Email) < 3 || err != nil {
		errList = append(errList, 6)
	} else if models.UserEmailIsExist(req.Email) {
		errList = append(errList, 7)
	}

	if len(errList) > 0 {
		returnJSON := make(map[string]interface{})
		returnJSON["id"] = 0
		returnJSON["err_list"] = errList
		c.JSON(http.StatusOK, returnJSON)
		return false
	}

	return true
}

// register (POST)
func Register(c *gin.Context) {
	var inputJSON SignupInfo

	if !registerReqDeal(c, &inputJSON) {
		return
	}

	var requireEmailActivate bool
	requireEmailActivate, _ = strconv.ParseBool(
		models.SvrSettingGet(svrsn.RequireEmailActivate))
	id, actToken, errList := models.UserCreate(inputJSON.LoginName, inputJSON.ShowName,
		inputJSON.Email, inputJSON.Password, 3, requireEmailActivate)

	if requireEmailActivate && len(errList) == 0 {
		var senderEmailAdrFormat = `"Simple Image Hosting" <` + models.SvrSettingGet(svrsn.SenderEmailAddress) + `>`
		var emailTitle = "Simple Image Hosting 註冊認證"
		var activateLink = "https://" + models.SvrSettingGet(svrsn.Hostname) + "/account-activate/" + actToken
		var bodyText = `
		<html>
  		<head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		</head>
  		<body>
		我們接受到了您在 ` + models.SvrSettingGet(svrsn.Hostname) + ` 的註冊申請<br>
		如果這不是您發出的請求，請忽略此訊息<br>
		<br>
		要完成認證請點擊下方連結:<br>
		<a href="` + activateLink + `" target="_blank">` + activateLink + `</a><br>
		</body>`

		err := common.GomailSender(
			senderEmailAdrFormat, inputJSON.Email, emailTitle, bodyText,
			models.SvrSettingGet(svrsn.SenderEmailServer),
			models.SvrSettingGet(svrsn.SenderEmailUser),
			models.SvrSettingGet(svrsn.SenderEmailPassword))
		if err != nil {
			errList = append(errList, 8)
			fmt.Println(err.Error())
			err = models.NotActUserDeleteByLoginName(inputJSON.LoginName)
			if err != nil {
				errList = append(errList, 8)
			}
		}
	}

	returnJSON := make(map[string]interface{})
	returnJSON["id"] = id
	returnJSON["err_list"] = errList

	c.JSON(http.StatusOK, returnJSON)
}

// AccountActivate (GET)
func AccountActivate(c *gin.Context) {
	actToken := c.Param("token")

	notActUser, exist, err := models.NotActUserGetByToken(actToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "Toekn is not valid or Expired."})
		return
	}

	notActUser.ID = 0
	var newUserID uint64
	newUserID, _ = models.ActivateUser(notActUser)
	if newUserID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error."})
		return
	}

	// 儲存使用者至 session
	session := sessions.Default(c)
	session.Set(userkey, newUserID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": newUserID, "show_name": notActUser.ShowName, "grade": notActUser.Grade, "msg": "Successfully activated."})
}

// Signin (POST)
func Signin(c *gin.Context) {
	var err error
	var inputJSON LoginInfo
	err = c.BindJSON(&inputJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request data"})
	}

	session := sessions.Default(c)

	// 檢查參數是否為空
	if inputJSON.LoginName == "" || inputJSON.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	var loginUser *models.User
	loginUser, err = models.UserGetByLoginName(inputJSON.LoginName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// 檢查使用者名稱和密碼是否正確
	err = bcrypt.CompareHashAndPassword(loginUser.PwdHash, []byte(inputJSON.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	_ = loginUser.UpdateLoginTime()

	// 儲存使用者至 session
	session.Set(userkey, loginUser.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": loginUser.ID, "show_name": loginUser.ShowName})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func Myinfo(c *gin.Context) {
	session := sessions.Default(c)
	idInterface := session.Get(userkey)
	if idInterface == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	id := idInterface.(uint64)

	user, userExist, err := models.UserGetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	if !userExist {
		session.Delete(userkey)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.String(http.StatusUnauthorized, "Unauthorized.")
		return
	}

	var avatarURL string
	if user.Avatar == "" {
		avatarURL = "/Avatars/default.png"
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "show_name": user.ShowName, "avatar": avatarURL, "grade": user.Grade})
}

func GetUserInfo(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request ID.")
		return
	}

	user, userExist, err := models.UserGetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	if !userExist {
		c.String(http.StatusNotFound, "User not found.")
		return
	}

	var avatarURL string
	if user.Avatar == "" {
		avatarURL = "/Avatars/default.png"
	}
	c.JSON(http.StatusOK, gin.H{"show_name": user.ShowName, "avatar": avatarURL, "introduction": user.Introduction})
}
