package models

import (
	"sih/models/svrsn"
	"strconv"
	"sync"

	"go.uber.org/zap"
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
		curV, exist, err := SvrSettingGetWithErr(name)
		if err != nil {
			return err
		}
		if !exist {
			err = SvrSettingCreate(name, defaultV)
			if err != nil {
				return err
			}
			curV = defaultV
		}
		svrSettingCache.Store(name, curV)
	}

	return nil
}

// SvrSetting server setting
type SvrSetting struct {
	Name  string `gorm:"primary_key"`
	Value string `gorm:"type:TEXT"`
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
	val, _, _ := SvrSettingGetWithErr(name)
	return val
}

// SvrSettingGet 藉由 name 獲取其 value
func SvrSettingGetWithErr(name svrsn.SvrSettingName) (val string, exist bool, err error) {
	// 先從 cache 讀取
	if cacheVal, cacheExist := svrSettingCache.Load(name); cacheExist {
		return cacheVal.(string), true, nil
	}

	// cache 中找不到再從 DB 讀取
	var setting = SvrSetting{Name: string(name)}
	res := db.Find(&setting)
	if res.Error == nil {
		if res.RowsAffected > 0 {
			val = setting.Value
			exist = true
			svrSettingCache.Store(name, setting.Value)
		}
	} else {
		logger.Error("find setting failed", zap.Error(res.Error))
		err = res.Error
	}

	return
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
