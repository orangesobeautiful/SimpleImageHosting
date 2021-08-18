package databaseoperation

import (
	"time"

	"gorm.io/gorm"
)

func CreateImage(title string, description string, imageType string, width int, height int, OwnerId int64) (int64, string, error) {
	//把資料寫入資料庫
	var newImage Image
	newImage.ID = 0
	newImage.Title = title
	newImage.Description = description
	newImage.Type = imageType
	newImage.Width = width
	newImage.Height = height
	newImage.OwnerID = OwnerId
	newImage.UpdateAt = time.Now().Unix()
	newImage.CreateAt = time.Now().Unix()

	res := db.Create(&newImage)
	if res.Error != nil {
		return 0, "", res.Error
	}

	//根據新資料的 ID 更新 HashID
	hashID, err := hashID.EncodeInt64([]int64{newImage.ID})
	if err != nil {
		return 0, "", err
	}

	var fileName string = hashID + "." + imageType
	res = db.Model(&newImage).Updates(map[string]interface{}{"HashID": hashID, "FileName": fileName})
	return newImage.ID, hashID, res.Error
}

func GetImageByID(id int64) (Image, *gorm.DB) {
	var image Image
	res := db.Where(&Image{ID: id}, "ID").Find(&image)
	return image, res
}

func GetImageListByOwnerID(ownerID int64) ([]Image, *gorm.DB) {
	var imageList []Image
	res := db.Where(&Image{OwnerID: ownerID}, "OwnerID").Find(&imageList)
	return imageList, res
}

func UpdateImage(id int64, updateData map[string]interface{}) *gorm.DB {
	var image Image = Image{ID: id}
	updateData["UpdateAt"] = time.Now().Unix()
	res := db.Model(&image).Updates(updateData)
	return res
}

func DeleteImage(id int64) *gorm.DB {
	var image Image = Image{ID: id}
	res := db.Where(&Image{ID: id}, "ID").Delete(&image)
	return res
}
