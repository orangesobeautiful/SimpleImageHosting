package databaseoperation

import (
	"crypto/rand"
	"encoding/base64"
	"net/mail"
	"time"
	"unicode/utf8"

	"SimpleImageHosting/common"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUser 新增使用者
func CreateUser(loginName string, showName string, email string, password string, grade int, requireEmailActivate bool) (int64, string, []int) {
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
	var actToken string

	if IsLoginNameUsed(loginName) {
		errList = append(errList, 1)
	}

	var loginNameLen int = len(loginName)
	if loginNameLen < 4 || loginNameLen > 30 {
		errList = append(errList, 2)
	}

	var showNameLen int = utf8.RuneCountInString(showName)
	if showNameLen < 1 || showNameLen > 15 {
		errList = append(errList, 3)
	}

	if len(password) < 6 {
		errList = append(errList, 4)
	}
	if len(email) > 256 {
		errList = append(errList, 5)
	}

	_, err := mail.ParseAddress(email)
	if len(email) < 3 || err != nil {
		errList = append(errList, 6)
	} else if IsEmailUsed(email) {
		errList = append(errList, 7)
	}

	var newUserID int64 = 0
	if len(errList) == 0 {
		bHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			errList = append(errList, 8)
			return newUserID, actToken, errList
		}

		var res *gorm.DB
		var currentTime int64 = time.Now().Unix()
		var newUser = User{LoginName: loginName, ShowName: showName, Email: email, PwdHash: bHash, Grade: grade, LastLoginTime: currentTime, CreatedAt: currentTime}
		if requireEmailActivate {
			var newNotActUser NotActivatedUser
			newNotActUser.User = newUser
			newNotActUser.EmailExpiration = common.MaxUnixTimeInt64

			// 在未被認證的使用者中 email 應該是要能重複的
			// 然而 Email 欄位是唯一限定不能重複(宣告時的繼承問題引起)
			// 因此在 NotActivatedUser 表格中使用同樣唯一限定的 loginName 暫時代替
			// 使用額外欄位 NotActEmail 紀錄真正的 Email
			newNotActUser.Email = loginName
			newNotActUser.NotActEmail = email

			randBytes := make([]byte, emailActTokenRandLen)
			if _, err := rand.Read(randBytes); err != nil {
				errList = append(errList, 8)
			}
			actToken = base64.RawURLEncoding.EncodeToString(randBytes)
			newNotActUser.ActaivateToken = actToken

			res = db.Create(&newNotActUser)
		} else {
			res = db.Create(&newUser)
		}
		if res.Error == nil {
			//newUser, _ = GetUserByLoginName(loginName)
			newUserID = newUser.ID
		} else {
			errList = append(errList, 8)
		}

	}

	return newUserID, actToken, errList
}

// IsLoginNameUsed 查詢 LoginName 是否被使用過
func IsLoginNameUsed(loginName string) bool {
	var user User
	res := db.Where(&User{LoginName: loginName}).Limit(1).Find(&user)
	if res.RowsAffected > 0 {
		return true
	}

	var notActUser NotActivatedUser
	notActUser.LoginName = loginName
	res = db.Where(&notActUser).Limit(1).Find(&notActUser)
	if res.RowsAffected > 0 {
		return true
	}
	return false
}

// IsEmailUsed 查詢 Email 是否被使用過
func IsEmailUsed(email string) bool {
	var user User
	//db.Where("Email = ?", email).First(&user)
	res := db.Where(&User{Email: email}).Limit(1).Find(&user)
	if res.RowsAffected > 0 {
		return true
	}
	return false
}

// IsUserExist 根據 ID 查詢 User 是否存在
func IsUserExist(userID int64) bool {
	var user User
	res := db.Where(&User{ID: userID}, "ID").Limit(1).Find(&user)
	if res.RowsAffected > 0 {
		return true
	}
	return false
}

// GetUserByLoginName 根據 LoginName 查詢 User
func GetUserByLoginName(loginName string) (User, *gorm.DB) {
	var user User
	var res *gorm.DB = db.Where(&User{LoginName: loginName}).Limit(1).Find(&user)
	return user, res
}

// GetUserByID 根據 ID 查詢 User
func GetUserByID(userID int64) (User, *gorm.DB) {
	var user User
	var res *gorm.DB = db.Where(&User{ID: userID}).Limit(1).Find(&user)
	return user, res
}

// UpdateLoginTime 更新使用者最後登入時間
func UpdateLoginTime(userID int64) *gorm.DB {
	var user User = User{ID: userID}
	res := db.Model(&user).Where(&user).Update("LastLoginTime", time.Now().Unix())
	return res
}

// GetNotActUserByToken 根據 Activate token 查詢 Not Activated User
func GetNotActUserByToken(actToken string) (NotActivatedUser, *gorm.DB) {
	var notActUser NotActivatedUser = NotActivatedUser{ActaivateToken: actToken}
	res := db.Model(&notActUser).Where(&notActUser).Limit(1).Find(&notActUser)
	return notActUser, res
}

// DeleteNotActUserByToken 根據 Activate token 查詢 Not Activated User
func DeleteNotActUserByLoginName(loginName string) *gorm.DB {
	var notActUser NotActivatedUser
	notActUser.LoginName = loginName
	res := db.Where(&notActUser).Delete(&NotActivatedUser{})
	return res
}

func ActivateUser(notActUser NotActivatedUser) (int64, *gorm.DB) {
	// 遷移資料從 NotActivatedUser 到 User
	var res *gorm.DB
	notActUser.MailVaild = true
	var currentTime int64 = time.Now().Unix()
	notActUser.User.Email = notActUser.NotActEmail
	notActUser.CreatedAt = currentTime
	res = db.Create(&notActUser.User)
	if res.Error != nil {
		return 0, res
	}

	res = db.Where(&NotActivatedUser{NotActEmail: notActUser.NotActEmail}).Delete(&NotActivatedUser{})
	return notActUser.ID, res
}
