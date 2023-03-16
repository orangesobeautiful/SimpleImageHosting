package models

import (
	"time"

	"gorm.io/gorm"
)

// Image datebase image struct
type Image struct {
	ID       int64  `gorm:"column: ID; type: BIGINT UNSIGNED NOT NULL auto_increment; primary_key;" json:"id"`
	HashID   string `gorm:"column: HashID; type: VARCHAR(30);" json:"hash_id"`
	OwnerID  int64  `gorm:"column: OwnerID; type: BIGINT UNSIGNED NOT NULL; index:idx_owner_id;" json:"owner_id"`
	FileName string `gorm:"column: FileName; type: VARCHAR(40);" json:"file_name"`

	Type       string `gorm:"column: Type; type: VARCHAR(10) NOT NULL;" json:"type"`
	Width      int    `gorm:"column: Width; type:INTEGER UNSIGNED NOT NULL " json:"width"`
	Height     int    `gorm:"column: Height; type:INTEGER UNSIGNED NOT NULL " json:"height"`
	Size       int64  `gorm:"column: Size; type:BIGINT UNSIGNED" json:"size"`
	MediumSize int64  `gorm:"column: MediumSize; type:BIGINT UNSIGNED; default:0" json:"medium_size"`

	Title       string `gorm:"column: Title; type:VARCHAR(30) NOT NULL;" json:"title"`
	Description string `gorm:"column: Description; type:VARCHAR(1000) NOT NULL;" json:"description"`
	CreateAt    int64  `gorm:"column: CreatedAt; type:BIGINT UNSIGNED NOT NULL " json:"create_at"`
	UpdateAt    int64  `gorm:"column: UpdateAt; type:BIGINT UNSIGNED NOT NULL " json:"update_at"`
}

// TableName 指定 Image 表格的名稱
func (Image) TableName() string {
	return "images"
}

// CreateImage 創建一筆新的 Image 紀錄
// return id, hashID, error
func CreateImage(title, description, imageType string, width, height int,
	ownerID int64) (id int64, hashIDStr string, err error) {
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
	hashIDStr, err = hashID.EncodeInt64([]int64{newImage.ID})
	if err != nil {
		return 0, "", err
	}

	var fileName = hashIDStr + "." + imageType
	res = db.Model(&newImage).Updates(map[string]interface{}{"HashID": hashID, "FileName": fileName})
	return newImage.ID, hashIDStr, res.Error
}

// ImageFindByHashID 藉由 hashID 獲取 image 資料
func ImageFindByHashID(id string) (*Image, error) {
	arrayRes, _ := hashID.DecodeInt64WithError(id)
	if len(arrayRes) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	imgID := arrayRes[0]

	return imageFindByID(imgID)
}

// imageFindByID 藉由 id 獲取 image 資料
func imageFindByID(id int64) (*Image, error) {
	image := &Image{}
	res := db.Where(&Image{ID: id}, "ID").Find(image)
	return image, res.Error
}

// ImageListByOwnerID 藉由 userID 獲取其擁有的圖片資料
func ImageListByOwnerID(ownerID int64) ([]*Image, error) {
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
