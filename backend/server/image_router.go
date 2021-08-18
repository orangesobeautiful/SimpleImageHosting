package server

import (
	"SimpleImageHosting/databaseoperation"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

		file, err := fileHeader.Open()
		var headerBytes []byte
		headerBytes = make([]byte, 14)
		_, err = file.Read(headerBytes)
		if err != nil {
			fmt.Println(err)
			return
		}

		resType := http.DetectContentType(headerBytes)
		var fileType string
		switch resType {
		case "image/bmp":
			fileType = "bmp"
			break
		case "image/gif":
			fileType = "gif"
			break
		case "image/x-icon":
			fileType = "ico"
			break
		case "image/jpeg":
			fileType = "jpg"
			break
		case "image/png":
			fileType = "png"
			break
		case "image/webp":
			fileType = "webp"
			break
		case "application/octet-stream":
			fileType = ""
			break
		default:
			fileType = ""
			break
		}

		var allowTypeList []string = []string{"jpg", "png", "gif"}
		var isAllowType bool = false

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
			err = file.Close()
			if err != nil {
				errList = append(errList, [2]string{fileHeader.Filename, "伺服器內部錯誤"})
				fmt.Println("file.Close() 錯誤", err.Error())
				continue
			}

			// 在資料庫中新增一筆紀錄
			imgID, imgHashID, err := databaseoperation.CreateImage("", "", fileType, imgConf.Width, imgConf.Height, userIDInt)
			var savePath = path.Join(imageSaveDir, imgHashID+"."+fileType)
			if err != nil {
				errList = append(errList, [2]string{fileHeader.Filename, "server error"})
			} else {
				// 上傳圖片
				err = c.SaveUploadedFile(fileHeader, savePath)
				if err != nil {
					// 如果上傳錯誤 則刪除那筆紀錄
					databaseoperation.DeleteImage(imgID)
					errList = append(errList, [2]string{fileHeader.Filename, "server error"})
					continue
				}

				// 讀取檔案大小
				fi, err := os.Stat(savePath)
				if err != nil {
					// Could not obtain stat, handle error
				} else {
					// 更新檔案大小
					res := databaseoperation.UpdateImage(imgID, map[string]interface{}{"Size": fi.Size()})
					if res.Error != nil {
						// 如果更新錯誤則刪除紀錄
						databaseoperation.DeleteImage(imgID)
						errList = append(errList, [2]string{fileHeader.Filename, "server error"})
						continue
					}
				}
			}

		} else {
			errList = append(errList, [2]string{fileHeader.Filename, "not allow type"})
		}
	}

	c.JSON(http.StatusOK, errList)
}

func editImage(c *gin.Context) {
	imgHashIDStr := c.Param("hashID")
	arrayRes, _ := hashID.DecodeInt64WithError(imgHashIDStr)
	if len(arrayRes) < 0 {
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
	c.BindJSON(&dataMap)

	img, _ := databaseoperation.GetImageByID(imgID)
	if img.OwnerID != userIDInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission."})
		return
	}

	// 為了避免用戶端上傳的 json 包含不允許修改的 key
	// 手動複製一份 map 只保留允許修改的項目
	var copyMap map[string]interface{} = make(map[string]interface{})
	for key, value := range dataMap {
		switch key {
		case "title":
			copyMap["Title"] = value
		case "description":
			copyMap["Description"] = value
			break
		default:
			break
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
	if len(arrayRes) < 0 {
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

	var avatarUrl string
	if user.Avatar == "" {
		avatarUrl = "/Avatars/default.png"
	}

	c.JSON(http.StatusOK, gin.H{
		"title":        img.Title,
		"owner_id":     user.ID,
		"owner_name":   user.ShowName,
		"owner_avatar": avatarUrl,
		"description":  img.Description,
		"original_url": imageDirectory + "/" + img.FileName,
		"create_at":    strconv.FormatInt(img.CreateAt, 10),
	})
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
		resList = append(resList, gin.H{
			"hash_id":      img.HashID,
			"title":        img.Title,
			"description":  img.Description,
			"original_url": imageDirectory + "/" + img.FileName,
		})
	}
	c.JSON(http.StatusOK, resList)
}
