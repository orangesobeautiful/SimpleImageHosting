package databaseoperation

import (
	"time"

	"gorm.io/gorm"
)

// CreateImage 創建一筆新的 Image 紀錄
// return id, hashID, error
func CreateImage(title string, description string, imageType string, width int, height int, OwnerID int64) (int64, string, error) {
	//把資料寫入資料庫
	var newImage Image
	newImage.ID = 0
	newImage.Title = title
	newImage.Description = description
	newImage.Type = imageType
	newImage.Width = width
	newImage.Height = height
	newImage.OwnerID = OwnerID
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

// GetImageByID 藉由 id 獲取 image 資料
func GetImageByID(id int64) (Image, *gorm.DB) {
	var image Image
	res := db.Where(&Image{ID: id}, "ID").Find(&image)
	return image, res
}

// GetImageListByOwnerID 藉由 userID 獲取其擁有的圖片資料
func GetImageListByOwnerID(ownerID int64) ([]Image, *gorm.DB) {
	var imageList []Image
	res := db.Where(&Image{OwnerID: ownerID}, "OwnerID").Find(&imageList)
	return imageList, res
}

// UpdateImage 更新圖片資料，欄位由 updateData 的 key value 決定
func UpdateImage(id int64, updateData map[string]interface{}) *gorm.DB {
	var image = Image{ID: id}
	updateData["UpdateAt"] = time.Now().Unix()
	res := db.Model(&image).Updates(updateData)
	return res
}

// DeleteImage 刪除圖片紀錄
func DeleteImage(id int64) *gorm.DB {
	var image = Image{ID: id}
	res := db.Where(&Image{ID: id}, "ID").Delete(&image)
	return res
}
