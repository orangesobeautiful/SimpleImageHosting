package databaseoperation

import "gorm.io/gorm"

func CreateSetting(name string, value string) *gorm.DB {
	var newSetting Setting = Setting{Name: name, Value: value}
	res := db.Create(&newSetting)
	return res
}

func GetSetting(name string) (string, *gorm.DB) {
	var setting Setting = Setting{Name: name}
	res := db.First(&setting)
	return setting.Value, res

}

func UpdateSetting(name string, value string) *gorm.DB {
	var setting Setting = Setting{Name: name, Value: value}
	res := db.Model(&setting).Where(&Setting{Name: name}).Update("Value", value)
	return res
}
