package databaseoperation

import "gorm.io/gorm"

// CreateSetting 創建一筆新的 Setting 紀錄
func CreateSetting(name string, value string) *gorm.DB {
	var newSetting = Setting{Name: name, Value: value}
	res := db.Create(&newSetting)
	return res
}

// GetSetting 藉由 name 獲取其 value
func GetSetting(name string) (string, *gorm.DB) {
	var setting = Setting{Name: name}
	res := db.First(&setting)
	return setting.Value, res

}

// UpdateSetting 更新設定資料
func UpdateSetting(name string, value string) *gorm.DB {
	var setting = Setting{Name: name, Value: value}
	res := db.Model(&setting).Where(&Setting{Name: name}).Update("Value", value)
	return res
}
