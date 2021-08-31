package server

import (
	"SimpleImageHosting/common"
	"SimpleImageHosting/databaseoperation"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getServerInfo(c *gin.Context) {
	var requireEmailActivate bool
	requireEmailActivate, _ = strconv.ParseBool(webSetting["RequireEmailActivate"])
	c.JSON(http.StatusOK, gin.H{
		"require_email_activate": requireEmailActivate,
	})
}

func getWebsiteSettings(c *gin.Context) {
	var res *gorm.DB
	var err error

	// 驗證使用者身分
	session := sessions.Default(c)
	idInterface := session.Get(userkey)
	if idInterface == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	userID := idInterface.(int64)

	user, res := databaseoperation.GetUserByID(userID)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else if res.RowsAffected <= 0 {
		session.Delete(userkey)
		if err = session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.String(http.StatusUnauthorized, "Unauthorized.")
		return
	}
	if user.Grade != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission."})
		return
	}

	returnJSON := make(map[string]interface{})

	returnJSON["require_email_activate"], _ = strconv.ParseBool(webSetting["RequireEmailActivate"])
	returnJSON["sender_email_server"] = webSetting["SenderEmailServer"]
	returnJSON["sender_email_address"] = webSetting["SenderEmailAddress"]
	returnJSON["sender_email_user"] = webSetting["SenderEmailUser"]
	returnJSON["hostname"] = webSetting["Hostname"]

	c.JSON(http.StatusOK, returnJSON)
}

// websiteSetting 修改網站設定
func editWebsiteSettings(c *gin.Context) {
	session := sessions.Default(c)
	idInterface := session.Get(userkey)
	if idInterface == nil {
		// 錯誤的cookie
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	id := idInterface.(int64)

	var user databaseoperation.User
	var res *gorm.DB
	var err error
	user, res = databaseoperation.GetUserByID(id)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else if res.RowsAffected <= 0 {
		// 資料庫裡找不到使用者
		session.Delete(userkey)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized."})
		return
	}

	if user.Grade != 1 {
		// 權限不符
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission."})
		return
	}

	// 取得傳輸的JSON
	var dataMap map[string]interface{}
	err = c.BindJSON(&dataMap)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request data"})
	}

	// 為了避免用戶端上傳的 json 包含不允許修改的 key
	// 手動複製一份 map 只保留允許修改的項目
	var copyMap = make(map[string]interface{})
	for key, value := range dataMap {
		switch key {
		case "hostname":
			copyMap["Hostname"] = value
		case "require_email_activate":
			// 檢驗是否為正確 boolean 值
			switch checkV := value.(type) {
			case bool:
				copyMap["RequireEmailActivate"] = strconv.FormatBool(checkV)
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "The value of \"RequireEmailActivate\" is not boolean"})
				return
			}
		case "sender_email_server":
			copyMap["SenderEmailServer"] = value
		case "sender_email_address":
			copyMap["SenderEmailAddress"] = value
		case "sender_email_user":
			copyMap["SenderEmailUser"] = value
		case "sender_email_password":
			copyMap["SenderEmailPassword"] = value
		default:
		}
	}

	// 檢查 email 設定
	var senderEmailServer, senderEmailUser, senderEmailPassword string = "", "", ""
	var emailEmptyCount = 0
	if val, ok := copyMap["SenderEmailServer"]; ok && val.(string) != "" {
		senderEmailServer = val.(string)
	} else {
		emailEmptyCount++
	}
	if val, ok := copyMap["SenderEmailAddress"]; ok && val.(string) != "" {
		_ = val.(string)
	} else {
		emailEmptyCount++
	}
	if val, ok := copyMap["SenderEmailUser"]; ok && val.(string) != "" {
		senderEmailUser = val.(string)
	} else {
		emailEmptyCount++
	}
	if val, ok := copyMap["SenderEmailPassword"]; ok && val.(string) != "" {
		senderEmailPassword = val.(string)
	} else {
		emailEmptyCount++
	}

	fmt.Println(emailEmptyCount)
	switch emailEmptyCount {
	case 0:
		err = common.CheckEmailLogin(senderEmailServer, senderEmailUser, senderEmailPassword)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	case 1, 2, 3:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email server, address, user or password is empty."})
		return
	case 4:
		// 如果設定中存在 RequireEmailActivate 且為 true
		if val, ok := copyMap["RequireEmailActivate"]; ok && val.(string) == strconv.FormatBool(true) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "RequireEmailActivate has been set. But data does not contain email settings."})
			return
		}
	}

	// 寫入資料庫
	if val, ok := copyMap["RequireEmailActivate"]; ok {
		webSetting["RequireEmailActivate"] = val.(string)
	}
	for key, value := range copyMap {
		res = databaseoperation.UpdateSetting(key, value.(string))
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		webSetting[key] = value.(string)
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Update success."})
}
