package controller

import (
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"

	// user for image.DecodeConfig
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"sih/config"
	"sih/models"

	"github.com/disintegration/imaging"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const mdWidth = 700

// fileHeaderDeal 處理上傳的 fileHeader
// 發生錯誤時回傳 [2]string{"Filename", "錯誤訊息"}
// 順利執行完時回傳空值皆為字串
func imgFHeaderDeal(userIDInt int64, fileHeader *multipart.FileHeader) [2]string {
	var imgID int64
	var needRecycle = false
	var gnrDBRecord, gnrMDImageFile, gnrOrgImageFile bool = false, false, false
	var savePath, mdSavePath string
	file, err := fileHeader.Open()
	if err != nil {
		return [2]string{fileHeader.Filename, "伺服器內部錯誤"}
	}
	defer file.Close()
	var headerBytes = make([]byte, 14)
	_, err = file.Read(headerBytes)
	if err != nil {
		return [2]string{fileHeader.Filename, "伺服器內部錯誤"}
	}

	var resType = http.DetectContentType(headerBytes)
	var fileType string
	switch resType {
	case "image/bmp":
		fileType = "bmp"
	case "image/gif":
		fileType = "gif"
	case "image/x-icon":
		fileType = "ico"
	case "image/jpeg":
		fileType = "jpg"
	case "image/png":
		fileType = "png"
	case "image/webp":
		fileType = "webp"
	case "application/octet-stream":
		fileType = ""
	default:
		fileType = ""
	}

	var allowTypeList = []string{"jpg", "png", "gif"}
	var isAllowType = false

	for _, allowType := range allowTypeList {
		if fileType == allowType {
			isAllowType = true
			break
		}
	}

	if isAllowType {
		//將 seek 重新移動到檔案開頭
		_, err = file.Seek(0, 0)
		if err != nil {
			return [2]string{fileHeader.Filename, "伺服器內部錯誤"}
		}
		imgConf, _, err := image.DecodeConfig(file)
		if err != nil {
			return [2]string{fileHeader.Filename, "無法解析圖片"}
		}

		// 在資料庫中新增一筆紀錄
		// (需要先新增紀錄產生圖片的 HashID 才有辦法確認要儲存的檔名)
		var mdOut, out *os.File
		var fi os.FileInfo
		var orgImgFileSize, mdImgFileSize int64
		var imgHashID string
		imgID, imgHashID, err = models.CreateImage("", "", fileType, imgConf.Width, imgConf.Height, userIDInt)
		savePath = path.Join(cfg.ImageMDSaveDir(), imgHashID+"."+fileType)
		mdSavePath = path.Join(cfg.ImageMDSaveDir(), config.ImageMDDirectory, imgHashID+".md."+fileType)
		if err != nil {
			return [2]string{fileHeader.Filename, "server error"}
		}

		//在生成資料庫紀錄、圖片檔案後發生錯誤時，刪除這些已經寫入的東西
		defer func() {
			if needRecycle {
				if gnrDBRecord {
					// TODO: 引用 create 後的 image
					delImg := &models.Image{ID: imgID}
					delImg.Delete()
				}
				if gnrMDImageFile {
					os.Remove(savePath)
				}
				if gnrOrgImageFile {
					os.Remove(mdSavePath)
				}
			}
		}()

		gnrDBRecord = true
		// 產生縮圖
		if imgConf.Width > mdWidth {
			_, err = file.Seek(0, 0)
			if err != nil {
				needRecycle = true
				return [2]string{fileHeader.Filename, "伺服器內部錯誤"}
			}
			var srcImg, resizeImg image.Image
			srcImg, err = imaging.Decode(file, imaging.AutoOrientation(true))
			if err != nil {
				needRecycle = true
				return [2]string{fileHeader.Filename, "無法讀取圖片"}
			}
			resizeImg = imaging.Resize(srcImg, mdWidth, 0, imaging.Linear)
			mdOut, err = os.Create(mdSavePath)
			if err != nil {
				needRecycle = true
				return [2]string{fileHeader.Filename, "伺服器內部錯誤"}
			}
			defer func() {
				if mdOut.Close() != nil {
					needRecycle = true
				}
			}()
			err = imaging.Encode(mdOut, resizeImg, imaging.JPEG)
			if err != nil {
				needRecycle = true
				return [2]string{fileHeader.Filename, "伺服器內部錯誤"}
			}
			gnrMDImageFile = true
			fi, err = os.Stat(mdSavePath)
			if err != nil {
				needRecycle = true
				return [2]string{fileHeader.Filename, "伺服器內部錯誤"}
			}
			mdImgFileSize = fi.Size()
		}

		// 儲存原圖
		_, err = file.Seek(0, 0)
		if err != nil {
			needRecycle = true
			return [2]string{fileHeader.Filename, "伺服器內部錯誤"}
		}
		out, err = os.Create(savePath)
		if err != nil {
			needRecycle = true
			return [2]string{fileHeader.Filename, "伺服器內部錯誤"}
		}
		defer func() {
			if out.Close() != nil {
				needRecycle = true
			}
		}()

		_, err = io.Copy(out, file)
		if err != nil {
			needRecycle = true
			return [2]string{fileHeader.Filename, "伺服器內部錯誤"}
		}
		gnrOrgImageFile = true

		// 讀取檔案大小
		fi, err = os.Stat(savePath)
		if err != nil {
			needRecycle = true
			return [2]string{fileHeader.Filename, "伺服器內部錯誤"}
		}
		orgImgFileSize = fi.Size()
		// 更新檔案大小
		// TODO: 引用 create 後的 image
		updateImg := &models.Image{ID: imgID}
		err = updateImg.Update(map[string]interface{}{"Size": orgImgFileSize, "MediumSize": mdImgFileSize})
		if err != nil {
			needRecycle = true
			return [2]string{fileHeader.Filename, "伺服器內部錯誤"}
		}
	} else {
		return [2]string{fileHeader.Filename, "not allow type"}
	}
	return [2]string{"", ""}
}

func UploadImage(c *gin.Context) {
	session := sessions.Default(c)
	form, err := c.MultipartForm()
	//fileHeaders 的讀取需要放在錯誤偵測的 return 前面
	//如果在讀取前就回傳，會被當作伺服端莫名中斷連線，進而回傳伺服端錯誤的代碼。
	//例如在nginx中會回傳 502
	fileHeaders := form.File["image[]"]

	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you have to sign in."})
		return
	}

	userIDInt, _ := user.(int64)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request."})
		return
	}

	if len(fileHeaders) == 0 {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Upload Images."})
		return
	}

	//errList (第0項: file name，第1項: 錯誤原因)
	var errList [][2]string

	// 處理上傳的檔案
	for _, fileHeader := range fileHeaders {
		var errMsg = imgFHeaderDeal(userIDInt, fileHeader)
		if errMsg[0] != "" {
			errList = append(errList, errMsg)
		}
	}

	c.JSON(http.StatusOK, errList)
}

func EditImage(c *gin.Context) {
	var err error
	imgHashIDStr := c.Param("hashID")

	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you have to sign in."})
		return
	}
	userIDInt, _ := user.(int64)

	var dataMap map[string]interface{}
	err = c.BindJSON(&dataMap)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request data"})
	}

	img, err := models.ImageFindByHashID(imgHashIDStr)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.String(http.StatusNotFound, "image not found")
		} else {
			c.String(http.StatusInternalServerError, "internal server error")
		}
	}
	if img.OwnerID != userIDInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission."})
		return
	}

	// 為了避免用戶端上傳的 json 包含不允許修改的 key
	// 手動複製一份 map 只保留允許修改的項目
	var copyMap = make(map[string]interface{})
	for key, value := range dataMap {
		switch key {
		case "title":
			copyMap["Title"] = value
		case "description":
			copyMap["Description"] = value
		default:
		}
	}

	err = img.Update(copyMap)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.String(http.StatusOK, "success update.")
}

func DeleteImage(c *gin.Context) {
	imgHashIDStr := c.Param("hashID")

	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you have to sign in."})
		return
	}
	userIDInt, _ := user.(int64)

	img, err := models.ImageFindByHashID(imgHashIDStr)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.String(http.StatusNotFound, "image not found")
		} else {
			c.String(http.StatusInternalServerError, "internal server error")
		}
	}
	var filename = img.FileName
	if img.OwnerID != userIDInt {
		c.String(http.StatusForbidden, "forbidden.")
		return
	}

	err = img.Delete()
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = os.Remove(path.Join(cfg.ImageSaveDir(), filename))
	if err != nil {
		fmt.Println(err)
	}

	c.String(http.StatusOK, "success delete.")
}

func GetImage(c *gin.Context) {
	imgHashIDStr := c.Param("hashID")

	img, err := models.ImageFindByHashID(imgHashIDStr)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.String(http.StatusNotFound, "image not found")
		} else {
			c.String(http.StatusInternalServerError, "internal server error")
		}
	}

	user, res := models.GetUserByID(img.OwnerID)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else if res.RowsAffected <= 0 {
		c.String(http.StatusNotFound, "Owner not found.")
		return
	}

	var avatarURL string
	if user.Avatar == "" {
		avatarURL = "/Avatars/default.png"
	} else {
		avatarURL = user.Avatar
	}

	var imgDataJSON = gin.H{
		"title":        img.Title,
		"owner_id":     user.ID,
		"owner_name":   user.ShowName,
		"owner_avatar": avatarURL,
		"description":  img.Description,
		"original_url": path.Join(config.ImageDirectory, img.HashID+"."+img.Type),
		"create_at":    strconv.FormatInt(img.CreateAt, 10),
	}
	if img.MediumSize > 0 {
		imgDataJSON["md_url"] = path.Join(config.ImageDirectory, config.ImageMDDirectory, img.HashID+".md."+img.Type)
	}
	c.JSON(http.StatusOK, imgDataJSON)
}

func GetUserImages(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request ID.")
		return
	}

	if !models.IsUserExist(userID) {
		c.String(http.StatusNotFound, "User does not exist.")
		return
	}

	imageList, err := models.ImageListByOwnerID(userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error.")
		return
	}

	var resList []gin.H
	for _, img := range imageList {
		var imgDataJSON = gin.H{
			"hash_id":      img.HashID,
			"title":        img.Title,
			"description":  img.Description,
			"original_url": path.Join(config.ImageDirectory, img.HashID+"."+img.Type),
		}
		if img.MediumSize > 0 {
			imgDataJSON["md_url"] = path.Join(config.ImageDirectory, config.ImageMDDirectory, img.HashID+".md."+img.Type)
		}
		resList = append(resList, imgDataJSON)
	}
	c.JSON(http.StatusOK, resList)
}
