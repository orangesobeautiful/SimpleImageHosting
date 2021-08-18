package databaseoperation

import (
	"strconv"

	"github.com/speps/go-hashids/v2"
	"gorm.io/gorm"
)

var db *gorm.DB
var hashID *hashids.HashID

var emailActTokenRandLen int = 32

// SetDB set databae
func SetHashIDData(salt string, minLen int) error {
	hashIDData := hashids.NewData()
	hashIDData.Salt = salt
	hashIDData.MinLength = minLen
	var err error
	hashID, err = hashids.NewWithData(hashIDData)
	return err
}

// SetDB set databae
func SetDB(inputdb *gorm.DB) {
	db = inputdb
}

// InitDatabase 初始化資料庫
func InitDatabase() {
	db.AutoMigrate(Setting{})
	db.AutoMigrate(User{})
	db.AutoMigrate(Image{})
	db.AutoMigrate(NotActivatedUser{})

	//讀取 setting 列表 (第0項: setting name，第1項: 讀取不到時創建的預設值)
	settingTable := [][2]string{{"SessionSecretKey", ""},
		{"OwnerRegistered", strconv.FormatBool(false)},
		{"HashIDSalt", ""},
		{"Hostname", ""},
		{"RequireEmailActivate", strconv.FormatBool(false)},
		{"SenderEmailServer", ""},
		{"SenderEmailAddress", ""},
		{"SenderEmailUser", ""},
		{"SenderEmailPassword", ""}}
	for _, v := range settingTable {
		if _, res := GetSetting(v[0]); res.Error != nil {
			CreateSetting(v[0], v[1])
		}
	}

}
