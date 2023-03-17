package models

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"sih/common"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User datebase user struct
type User struct {
	ID        uint64 `gorm:"auto_increment; primary_key"`
	LoginName string `gorm:"type:VARCHAR(30); uniqueIndex"`
	ShowName  string `gorm:"type:VARCHAR(30)"`
	Email     string `gorm:"type:VARCHAR(255); uniqueIndex"`

	PwdHash      []byte `gorm:"type:BINARY(128)" json:"-"`
	Avatar       string `gorm:"type:VARCHAR(30)"`
	Introduction string `gorm:"type:VARCHAR(500)"`
	Grade        int
	MailVaild    bool
	CreatedAt    int64 `gorm:"index"`
	UpdateAt     int64
}

// TableName 指定 User 表格的名稱
func (User) TableName() string {
	return "users"
}

// NotActivatedUser 未進行郵件認證的使用者
// 真正紀錄 Email 的欄位是 NotActEmail
type NotActivatedUser struct {
	User
	NotActEmail     string `gorm:"type:VARCHAR(255); index"`
	ActaivateToken  string `gorm:"type:VARCHAR(344); uniqueIndex"`
	EmailExpiration int64  `gorm:"type:BIGINT UNSIGNED"`
}

// TableName 指定 NotActivatedUser 表格的名稱
func (NotActivatedUser) TableName() string {
	return "not_activated_user"
}

// UserCreate 新增使用者
func UserCreate(loginName, showName, email, password string,
	grade int, requireEmailActivate bool) (id uint64, actToken string, errList []int) {
	var newUserID uint64
	const unknowErrCode = 8

	bHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errList = append(errList, unknowErrCode)
		return newUserID, actToken, errList
	}
	var res *gorm.DB
	var currentTime = time.Now().Unix()
	var newUser = User{
		LoginName: loginName,
		ShowName:  showName,
		Email:     email,
		PwdHash:   bHash,
		Grade:     grade,
		CreatedAt: currentTime,
		UpdateAt:  currentTime,
	}
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

		const emailActTokenRandLen = 32
		randBytes := make([]byte, emailActTokenRandLen)
		if _, err := rand.Read(randBytes); err != nil {
			errList = append(errList, unknowErrCode)
		}
		actToken = base64.RawURLEncoding.EncodeToString(randBytes)
		newNotActUser.ActaivateToken = actToken

		res = db.Create(&newNotActUser)
	} else {
		res = db.Create(&newUser)
	}
	if res.Error == nil {
		newUserID = newUser.ID
	} else {
		errList = append(errList, unknowErrCode)
	}

	return newUserID, actToken, errList
}

// IsLoginNameUsed 查詢 LoginName 是否被使用過
func IsLoginNameUsed(loginName string) bool {
	var user User
	res := db.Where(&User{LoginName: loginName}).First(&user)
	if res.RowsAffected > 0 {
		return true
	}

	var notActUser NotActivatedUser
	notActUser.LoginName = loginName
	res = db.Where(&notActUser).First(&notActUser)
	return res.RowsAffected > 0
}

// UserEmailIsExist 查詢 Email 是否被使用過
func UserEmailIsExist(email string) bool {
	user := &User{Email: email}
	res := db.Where(user).Limit(1).First(&user)
	return res.RowsAffected > 0
}

// UserIsExist 根據 ID 查詢 User 是否存在
func UserIsExist(userID uint64) bool {
	var user User
	res := db.Where(&User{ID: userID}, "ID").First(&user)
	return res.RowsAffected > 0
}

// UserGetByLoginName 根據 LoginName 查詢 User
func UserGetByLoginName(loginName string) (*User, error) {
	user := &User{LoginName: loginName}
	var res = db.Where(user).First(user)
	return user, res.Error
}

// UserGetByID 根據 ID 查詢 User
func UserGetByID(userID uint64) (resUser *User, exist bool, err error) {
	user := &User{ID: userID}
	res := db.Find(&user)
	if res.Error == nil {
		if res.RowsAffected > 0 {
			resUser = user
			exist = true
		}
	} else {
		err = res.Error
	}
	return
}

// UserUpdateLoginTime 更新使用者最後登入時間
func (user *User) UpdateLoginTime() error {
	// TODO: 新增登入紀錄 table
	return nil
}

// NotActUserGetByToken 根據 Activate token 查詢 Not Activated User
func NotActUserGetByToken(actToken string) (resUser *NotActivatedUser, exist bool, err error) {
	notActUser := &NotActivatedUser{ActaivateToken: actToken}
	res := db.Model(notActUser).Where(notActUser).Find(notActUser)
	if res.Error == nil {
		if res.RowsAffected > 0 {
			resUser = notActUser
			exist = true
		}
	} else {
		err = res.Error
	}

	return
}

// NotActUserDeleteByLoginName 根據 loginName 刪除 NotActUser 紀錄
func NotActUserDeleteByLoginName(loginName string) error {
	filter := &NotActivatedUser{}
	filter.LoginName = loginName

	res := db.Where(filter).Delete(&NotActivatedUser{})
	return res.Error
}

// ActivateUser avtivate a not activated user
func ActivateUser(notActUser *NotActivatedUser) (uint64, error) {
	// 遷移資料從 NotActivatedUser 到 User
	var res *gorm.DB
	notActUser.MailVaild = true
	var currentTime = time.Now().Unix()
	notActUser.User.Email = notActUser.NotActEmail
	notActUser.CreatedAt = currentTime
	res = db.Create(&notActUser.User)
	if res.Error != nil {
		return 0, res.Error
	}

	res = db.Where(&NotActivatedUser{NotActEmail: notActUser.NotActEmail}).Delete(&NotActivatedUser{})
	return notActUser.ID, res.Error
}
