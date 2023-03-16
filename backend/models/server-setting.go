package models

import (
	"sih/models/svrsn"
	"strconv"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// svrSettingCache serverSettingCache
var svrSettingCache sync.Map

// initServerSetting 初始化伺服器設定值
func initServerSetting() (err error) {
	defaultSetting := map[svrsn.SvrSettingName]string{
		svrsn.SessSecretKey:        "",
		svrsn.OwnerRegistered:      strconv.FormatBool(false),
		svrsn.HashIDSalt:           "",
		svrsn.Hostname:             "",
		svrsn.RequireEmailActivate: strconv.FormatBool(false),
		svrsn.SenderEmailServer:    "",
		svrsn.SenderEmailAddress:   "",
		svrsn.SenderEmailUser:      "",
		svrsn.SenderEmailPassword:  "",
	}
	for name, defaultV := range defaultSetting {
		if curV, err := SvrSettingGetWithErr(name); err == nil {
			svrSettingCache.Store(name, curV)
		} else {
			if err == gorm.ErrRecordNotFound {
				err = SvrSettingCreate(name, defaultV)
				if err != nil {
					return err
				}
				svrSettingCache.Store(name, defaultV)
			} else {
				return err
			}
		}
	}

	return nil
}

// SvrSetting server setting
type SvrSetting struct {
	Name  string `gorm:"column: Name; type:VARCHAR(30) NOT NULL; primary_key;" json:"name"`
	Value string `gorm:"column: Value; type:TEXT NOT NULL;" json:"value"`
}

// TableName 指定 Setting 表格的名稱
func (SvrSetting) TableName() string {
	return "settings"
}

// SvrSettingCreate 創建一筆新的 Server Setting
func SvrSettingCreate(name svrsn.SvrSettingName, value string) error {
	var newSetting = SvrSetting{Name: string(name), Value: value}
	res := db.Create(&newSetting)
	if res.Error != nil {
		logger.Error("create setting failed", zap.Error(res.Error))
		return res.Error
	}
	svrSettingCache.Store(name, value)

	return nil
}

// SvrSettingGet 藉由 name 獲取其 value，忽略 error
func SvrSettingGet(name svrsn.SvrSettingName) string {
	val, _ := SvrSettingGetWithErr(name)
	return val
}

// SvrSettingGet 藉由 name 獲取其 value
func SvrSettingGetWithErr(name svrsn.SvrSettingName) (string, error) {
	// 先從 cache 讀取
	if cacheVal, cacheExist := svrSettingCache.Load(name); cacheExist {
		return cacheVal.(string), nil
	}

	// cache 中找不到再從 DB 讀取
	var setting = SvrSetting{Name: string(name)}
	res := db.First(&setting)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", nil
		}
		logger.Error("find setting failed", zap.Error(res.Error))
		return "", res.Error
	}
	svrSettingCache.Store(name, setting.Value)

	return setting.Value, nil
}

// SvrSettingUpdate 更新設定資料
func SvrSettingUpdate(name svrsn.SvrSettingName, value string) error {
	var setting = SvrSetting{Name: string(name), Value: value}
	res := db.Model(&setting).Where(&SvrSetting{Name: string(name)}).Update("Value", value)
	if res.Error != nil {
		logger.Error("update setting failed", zap.Error(res.Error))
		return res.Error
	}
	svrSettingCache.Store(name, value)

	return nil
}
