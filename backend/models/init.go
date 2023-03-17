package models

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"sih/models/svrsn"

	"github.com/speps/go-hashids/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var db *gorm.DB
var logger *zap.Logger

var hashID *hashids.HashID

// initHashID 設定 HashID， 需要在 initServerSetting 之後執行
func initHashID() error {
	var err error

	// 初始化 hashids
	hashIDData := hashids.NewData()
	// 從資料庫讀取 salt
	hashIDSlat, exist, err := SvrSettingGetWithErr(svrsn.HashIDSalt)
	if err != nil {
		return err
	}

	// 沒設定過則自動生成並儲存
	if !exist {
		err = SvrSettingCreate(svrsn.HashIDSalt, "")
		if err != nil {
			logger.Error("create setting failed", zap.Error(err))
			return err
		}
	}
	if hashIDSlat == "" {
		const hashIDSlatLen = 8
		byteArray := make([]byte, hashIDSlatLen)
		if _, readErr := io.ReadFull(rand.Reader, byteArray); readErr != nil {
			logger.Error("generate rand bytes failed", zap.Error(readErr))
			return readErr
		}
		hashIDSlat = base64.StdEncoding.EncodeToString(byteArray)
		err = SvrSettingUpdate(svrsn.HashIDSalt, hashIDSlat)
		if err != nil {
			logger.Error("update HashIDSalt failed", zap.Error(err))
			return err
		}
	}
	hashIDData.Salt = hashIDSlat
	const minLen = 5
	hashIDData.MinLength = minLen
	hashID, err = hashids.NewWithData(hashIDData)
	if err != nil {
		logger.Error("hashids.NewWithData failed", zap.Error(err))
		return err
	}

	return nil
}

// Init 初始化
func Init(inputdb *gorm.DB, inLogger *zap.Logger) (err error) {
	db = inputdb
	logger = inLogger

	logger.Info("init models")

	logger.Info("init database")
	err = initDatabase()
	if err != nil {
		return
	}

	logger.Info("init image hashID")
	err = initHashID()
	if err != nil {
		return
	}

	return
}

// InitDatabase 初始化資料庫
func initDatabase() error {
	var err error

	autoMigrateList := map[string]any{
		"setting":          SvrSetting{},
		"user":             User{},
		"notActivatedUser": NotActivatedUser{},
		"image":            Image{},
	}
	for key, table := range autoMigrateList {
		err = db.AutoMigrate(table)
		if err != nil {
			logger.Error(key+" AutoMigrate failed", zap.Error(err))
			return err
		}
	}

	err = initServerSetting()
	if err != nil {
		return err
	}

	return nil
}
