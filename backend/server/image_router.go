package server

import (
	"SimpleImageHosting/databaseoperation"
	"fmt"
	"image"

	// user for image.DecodeConfig
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/disintegration/imaging"
)

func uploadImage(c *gin.Context) {
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

	for _, fileHeader := range fileHeaders {
		var needRecycle = false
		var imgID int64
		var gnrDBRecord, gnrMDImageFile, gnrOrgImageFile bool = false, false, false
		var savePath, mdSavePath string
		file, err := fileHeader.Open()
		if err != nil {
			errList = append(errList, [2]string{fileHeader.Filename, "伺服器內部錯誤"})
			continue
		}
		var headerBytes = make([]byte, 14)
		_, err = file.Read(headerBytes)
		if err != nil {
			errList = append(errList, [2]string{fileHeader.Filename, "伺服器內部錯誤"})
			continue
		}

		resType := http.DetectContentType(headerBytes)
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
				errList = append(errList, [2]string{fileHeader.Filename, "伺服器內部錯誤"})
				fmt.Println("file.Seek(0, 0) 錯誤", err.Error())
				continue
			}
			imgConf, _, err := image.DecodeConfig(file)
			if err != nil {
				fmt.Println("無法解析圖片", err.Error())
				errList = append(errList, [2]string{fileHeader.Filename, "無法解析圖片"})
				continue
			}

			// 在資料庫中新增一筆紀錄
			// (需要先新增紀錄產生圖片的 HashID 才有辦法確認要儲存的檔名)
			var mdOut, out *os.File
			var fi os.FileInfo
			var orgImgFileSize, mdImgFileSize int64
			var imgHashID string
			imgID, imgHashID, err = databaseoperation.CreateImage("", "", fileType, imgConf.Width, imgConf.Height, userIDInt)
			savePath = path.Join(imageSaveDir, imgHashID+"."+fileType)
			mdSavePath = path.Join(imageSaveDir, imageMDDirectory, imgHashID+".md."+fileType)
			if err != nil {
				errList = append(errList, [2]string{fileHeader.Filename, "server error"})
			} else {
				gnrDBRecord = true
				//產生縮圖
				if imgConf.Width > mdWidth {
					_, err = file.Seek(0, 0)
					if err != nil {
						errList = append(errList, [2]string{fileHeader.Filename, "伺服器內部錯誤"})
						needRecycle = true
						goto RecycleIntermediate
					}
					var srcImg, resizeImg image.Image
					srcImg, err = imaging.Decode(file, imaging.AutoOrientation(true))
					if err != nil {
						errList = append(errList, [2]string{fileHeader.Filename, "無法讀取圖片"})
						needRecycle = true
						goto RecycleIntermediate
					}
					resizeImg = imaging.Resize(srcImg, mdWidth, 0, imaging.Linear)
					mdOut, err = os.Create(mdSavePath)
					if err != nil {
						errList = append(errList, [2]string{fileHeader.Filename, "伺服器內部錯誤"})
						needRecycle = true
						goto RecycleIntermediate
					}
					err = imaging.Encode(mdOut, resizeImg, imaging.JPEG)
					if err != nil {
						errList = append(errList, [2]string{fileHeader.Filename, "伺服器內部錯誤"})
						needRecycle = true
						goto RecycleIntermediate
					}
					mdOut.Close()
					gnrMDImageFile = true
					fi, err = os.Stat(mdSavePath)
					if err != nil {
						errList = append(errList, [2]string{fileHeader.Filename, "伺服器內部錯誤"})
						needRecycle = true
						goto RecycleIntermediate
					}
					mdImgFileSize = fi.Size()
				}

				//儲存原圖
				_, err = file.Seek(0, 0)
				if err != nil {
					errList = append(errList, [2]string{fileHeader.Filename, "伺服器內部錯誤"})
					needRecycle = true
					goto RecycleIntermediate
				}
				out, err = os.Create(savePath)
				if err != nil {
					errList = append(errList, [2]string{fileHeader.Filename, "伺服器內部錯誤"})
					needRecycle = true
					goto RecycleIntermediate
				}
				_, err = io.Copy(out, file)
				if err != nil {
					errList = append(errList, [2]string{fileHeader.Filename, "伺服器內部錯誤"})
					needRecycle = true
					goto RecycleIntermediate
				}

				file.Close()
				out.Close()
				gnrOrgImageFile = true

				// 讀取檔案大小
				fi, err = os.Stat(savePath)
				if err != nil {
					errList = append(errList, [2]string{fileHeader.Filename, "伺服器內部錯誤"})
					needRecycle = true
					goto RecycleIntermediate
				}
				orgImgFileSize = fi.Size()
				// 更新檔案大小
				res := databaseoperation.UpdateImage(imgID, map[string]interface{}{"Size": orgImgFileSize, "MediumSize": mdImgFileSize})
				if res.Error != nil {
					errList = append(errList, [2]string{fileHeader.Filename, "server error"})
					needRecycle = true
					goto RecycleIntermediate
				}
			}

		} else {
			errList = append(errList, [2]string{fileHeader.Filename, "not allow type"})
		}

		//在生成資料庫紀錄、圖片檔案後發生錯誤時，刪除這些已經寫入的東西
	RecycleIntermediate:
		if needRecycle {
			if gnrDBRecord {
				databaseoperation.DeleteImage(imgID)
			}
			if gnrMDImageFile {
				os.Remove(savePath)
			}
			if gnrOrgImageFile {
				os.Remove(mdSavePath)
			}
		}

	}

	c.JSON(http.StatusOK, errList)
}

func editImage(c *gin.Context) {
	var err error
	imgHashIDStr := c.Param("hashID")
	arrayRes, _ := hashID.DecodeInt64WithError(imgHashIDStr)
	if len(arrayRes) == 0 {
		c.String(http.StatusBadRequest, "Bad request ID.")
		return
	}
	imgID := arrayRes[0]

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

	img, _ := databaseoperation.GetImageByID(imgID)
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

	res := databaseoperation.UpdateImage(imgID, copyMap)
	if res.Error != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.String(http.StatusOK, "success update.")
}

func deleteImage(c *gin.Context) {
	imgHashIDStr := c.Param("hashID")
	arrayRes, _ := hashID.DecodeInt64WithError(imgHashIDStr)
	if len(arrayRes) == 0 {
		c.String(http.StatusBadRequest, "Bad request ID.")
		return
	}
	imgID := arrayRes[0]

	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you have to sign in."})
		return
	}
	userIDInt, _ := user.(int64)

	img, _ := databaseoperation.GetImageByID(imgID)
	var filename string = img.FileName
	if img.OwnerID != userIDInt {
		c.String(http.StatusForbidden, "forbidden.")
		return
	}

	res := databaseoperation.DeleteImage(imgID)
	if res.Error != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err := os.Remove(path.Join(imageSaveDir, filename))
	if err != nil {
		fmt.Println(err)
	}

	c.String(http.StatusOK, "success delete.")
}

func getImage(c *gin.Context) {
	imgHashIDStr := c.Param("hashID")
	arrayRes, err := hashID.DecodeInt64WithError(imgHashIDStr)
	if err != nil {
		c.String(http.StatusNotFound, "Image not found.")
		return
	}

	if len(arrayRes) <= 0 {
		c.String(http.StatusBadRequest, "Bad request ID.")
		return
	}
	imgID := arrayRes[0]

	img, res := databaseoperation.GetImageByID(imgID)
	if res.RowsAffected <= 0 {
		c.String(http.StatusNotFound, "Image not found.")
		return
	}

	user, res := databaseoperation.GetUserByID(img.OwnerID)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else if res.RowsAffected <= 0 {
		c.String(http.StatusNotFound, "Owner not found.")
		return
	}

	var avatarURL string
	if user.Avatar == "" {
		avatarURL = "/Avatars/default.png"
	}

	var imgDataJSON = gin.H{
		"title":        img.Title,
		"owner_id":     user.ID,
		"owner_name":   user.ShowName,
		"owner_avatar": avatarURL,
		"description":  img.Description,
		"original_url": path.Join(imageDirectory, img.HashID+"."+img.Type),
		"create_at":    strconv.FormatInt(img.CreateAt, 10),
	}
	if img.MediumSize > 0 {
		imgDataJSON["md_url"] = path.Join(imageDirectory, imageMDDirectory, img.HashID+".md."+img.Type)
	}
	c.JSON(http.StatusOK, imgDataJSON)
}

func getUserImages(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request ID.")
		return
	}

	if !databaseoperation.IsUserExist(userID) {
		c.String(http.StatusNotFound, "User does not exist.")
		return
	}

	images, res := databaseoperation.GetImageListByOwnerID(userID)
	if res.Error != nil {
		c.String(http.StatusInternalServerError, "Internal server error.")
		return
	}

	var resList []gin.H
	for _, img := range images {
		var imgDataJSON = gin.H{
			"hash_id":      img.HashID,
			"title":        img.Title,
			"description":  img.Description,
			"original_url": path.Join(imageDirectory, img.HashID+"."+img.Type),
		}
		if img.MediumSize > 0 {
			imgDataJSON["md_url"] = path.Join(imageDirectory, imageMDDirectory, img.HashID+".md."+img.Type)
		}
		resList = append(resList, imgDataJSON)
	}
	c.JSON(http.StatusOK, resList)
}
