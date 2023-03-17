package models

import (
	"math"
	"time"
)

// Image datebase image struct
type Image struct {
	ID       uint64 `gorm:"auto_increment; primary_key"`
	HashID   string `gorm:"type:VARCHAR(30)"`
	OwnerID  uint64 `gorm:"index:idx_owner_id"`
	FileName string `gorm:"type:VARCHAR(255)"`

	Type       string `gorm:"type:VARCHAR(30)"`
	Width      uint32
	Height     uint32
	Size       int64
	MediumSize int64

	Title       string `gorm:"type:VARCHAR(30)"`
	Description string `gorm:"type:VARCHAR(1000)"`
	CreateAt    int64  `gorm:"type:BIGINT UNSIGNED"`
	UpdateAt    int64  `gorm:"type:BIGINT UNSIGNED"`
}

// TableName 指定 Image 表格的名稱
func (Image) TableName() string {
	return "images"
}

// CreateImage 創建一筆新的 Image 紀錄
// return id, hashID, error
func CreateImage(title, description, imageType string, width, height uint32,
	ownerID uint64) (id uint64, hashIDStr string, err error) {
	// 把資料寫入資料庫

	var newImage Image
	newImage.ID = 0
	newImage.Title = title
	newImage.Description = description
	newImage.Type = imageType
	newImage.Width = width
	newImage.Height = height
	newImage.OwnerID = ownerID
	newImage.UpdateAt = time.Now().Unix()
	newImage.CreateAt = time.Now().Unix()

	res := db.Create(&newImage)
	if res.Error != nil {
		return 0, "", res.Error
	}

	// 根據新資料的 ID 更新 HashID
	var numbers []int64
	if newImage.ID > math.MaxInt64 {
		numbers = []int64{int64(newImage.ID % math.MaxInt64), int64(newImage.ID / math.MaxInt64)}
	} else {
		numbers = []int64{int64(newImage.ID)}
	}
	hashIDStr, err = hashID.EncodeInt64(numbers)
	if err != nil {
		return 0, "", err
	}

	var fileName = hashIDStr + "." + imageType
	res = db.Model(&newImage).Updates(map[string]interface{}{"HashID": hashIDStr, "FileName": fileName})
	return newImage.ID, hashIDStr, res.Error
}

// ImageGetByHashID 藉由 hashID 獲取 image 資料
func ImageGetByHashID(id string) (img *Image, exist bool, err error) {
	arrayRes, _ := hashID.DecodeInt64WithError(id)
	var imgID uint64
	switch len(arrayRes) {
	case 0:
		return nil, false, nil
	case 1:
		imgID = uint64(arrayRes[0])
	case 2:
		imgID = uint64(arrayRes[0]) + uint64(arrayRes[1])*math.MaxInt64
	}

	return imageGetByID(imgID)
}

// imageGetByID 藉由 id 獲取 image 資料
func imageGetByID(id uint64) (img *Image, exist bool, err error) {
	image := &Image{}
	res := db.Where(&Image{ID: id}).Find(image)
	if res.Error == nil {
		if res.RowsAffected > 0 {
			img = image
			exist = true
		}
	} else {
		err = res.Error
	}

	return
}

// ImageListByOwnerID 藉由 userID 獲取其擁有的圖片資料
func ImageListByOwnerID(ownerID uint64) ([]*Image, error) {
	var imageList []*Image
	res := db.Where(&Image{OwnerID: ownerID}, "OwnerID").Find(&imageList)
	return imageList, res.Error
}

// Update 更新圖片資料，欄位由 updateData 的 key value 決定
func (img *Image) Update(updateData map[string]interface{}) error {
	var image = Image{ID: img.ID}
	updateData["UpdateAt"] = time.Now().Unix()
	res := db.Model(image).Updates(updateData)
	return res.Error
}

// Delete 刪除圖片紀錄
func (img *Image) Delete() error {
	res := db.Where(&Image{ID: img.ID}, "ID").Delete(img)
	return res.Error
}
