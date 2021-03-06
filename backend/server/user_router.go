package server

import (
	"SimpleImageHosting/common"
	"SimpleImageHosting/databaseoperation"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
func siteOwnerRegister(c *gin.Context) {
	var err error
	if ownerRegistered {
		c.String(http.StatusNotFound, "404 page not found")
	} else {
		var inputJSON SignupInfo
		err = c.BindJSON(&inputJSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request data"})
		}

		id, _, errList := databaseoperation.CreateUser(inputJSON.LoginName, inputJSON.ShowName, inputJSON.Email, inputJSON.Password, 1, false)

		if len(errList) == 0 {
			databaseoperation.UpdateSetting("OwnerRegistered", strconv.FormatBool(true))
			ownerRegistered = true
		}

		res := make(map[string]interface{})
		res["id"] = id
		res["errList"] = errList

		c.JSON(http.StatusOK, res)
		return
	}
}

// register (POST)
func register(c *gin.Context) {
	var err error
	var res *gorm.DB
	var inputJSON SignupInfo
	err = c.BindJSON(&inputJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request data"})
	}

	var requireEmailActivate bool
	requireEmailActivate, _ = strconv.ParseBool(webSetting["RequireEmailActivate"])
	id, actToken, errList := databaseoperation.CreateUser(inputJSON.LoginName, inputJSON.ShowName, inputJSON.Email, inputJSON.Password, 3, requireEmailActivate)

	if requireEmailActivate && len(errList) == 0 {
		var senderEmailAdrFormat = `"Simple Image Hosting" <` + webSetting["SenderEmailAddress"] + `>`
		var emailTitle = "Simple Image Hosting 註冊認證"
		var activateLink = "https://" + webSetting["Hostname"] + "/account-activate/" + actToken
		var bodyText = `
		<html>
  		<head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		</head>
  		<body>
		我們接受到了您在 ` + webSetting["Hostname"] + ` 的註冊申請<br>
		如果這不是您發出的請求，請忽略此訊息<br>
		<br>
		要完成認證請點擊下方連結:<br>
		<a href="` + activateLink + `" target="_blank">` + activateLink + `</a><br>
		</body>`

		err := common.GomailSender(senderEmailAdrFormat, inputJSON.Email, emailTitle, bodyText, webSetting["SenderEmailServer"], webSetting["SenderEmailUser"], webSetting["SenderEmailPassword"])
		if err != nil {
			errList = append(errList, 8)
			fmt.Println(err.Error())
			res = databaseoperation.DeleteNotActUserByLoginName(inputJSON.LoginName)
			if res.Error != nil {
				errList = append(errList, 8)
			}
		}
	}

	returnJSON := make(map[string]interface{})
	returnJSON["id"] = id
	returnJSON["err_list"] = errList

	c.JSON(http.StatusOK, returnJSON)
}

// accountActivate (GET)
func accountActivate(c *gin.Context) {
	actToken := c.Param("token")
	var res *gorm.DB
	var notActUser databaseoperation.NotActivatedUser
	notActUser, res = databaseoperation.GetNotActUserByToken(actToken)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	} else if res.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Toekn is not valid or Expired."})
		return
	}

	notActUser.ID = 0
	var newUserID int64
	newUserID, _ = databaseoperation.ActivateUser(notActUser)
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

// signin (POST)
func signin(c *gin.Context) {
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

	var loginUser databaseoperation.User
	var res *gorm.DB
	loginUser, res = databaseoperation.GetUserByLoginName(inputJSON.LoginName)
	if res.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// 檢查使用者名稱和密碼是否正確
	err = bcrypt.CompareHashAndPassword(loginUser.PwdHash, []byte(inputJSON.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	res = databaseoperation.UpdateLoginTime(loginUser.ID)
	if res.Error != nil {
		fmt.Println("update user login time error:", err)
	}

	// 儲存使用者至 session
	session.Set(userkey, loginUser.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": loginUser.ID, "show_name": loginUser.ShowName})
}

func logout(c *gin.Context) {
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

func myinfo(c *gin.Context) {
	session := sessions.Default(c)
	idInterface := session.Get(userkey)
	if idInterface == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	id := idInterface.(int64)

	user, res := databaseoperation.GetUserByID(id)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else if res.RowsAffected <= 0 {
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

func getUserInfo(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request ID.")
		return
	}

	user, res := databaseoperation.GetUserByID(userID)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else if res.RowsAffected <= 0 {
		c.String(http.StatusNotFound, "User not found.")
		return
	}

	var avatarURL string
	if user.Avatar == "" {
		avatarURL = "/Avatars/default.png"
	}
	c.JSON(http.StatusOK, gin.H{"show_name": user.ShowName, "avatar": avatarURL, "introduction": user.Introduction})
}
