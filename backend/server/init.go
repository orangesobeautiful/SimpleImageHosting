package server

import (
	"SimpleImageHosting/common"
	"SimpleImageHosting/databaseoperation"
	"SimpleImageHosting/settingoperation"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/speps/go-hashids/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	userkey = "user"
)

var mdWidth = 700

var fileSaveDir string
var avatarDirectory = "Avatars"
var imageDirectory = "ImagesDirect"
var imageMDDirectory = "Medium"
var avatarSaveDir string
var imageSaveDir string

var hashID *hashids.HashID

var ownerRegistered bool

var webSetting map[string]string

// Start start server
func Start() {
	webSetting = make(map[string]string)
	var setting = readSetting()
	var ioWriterList []io.Writer
	var err error

	checkDictoryStruct()

	if setting.Server.Stdout {
		ioWriterList = append(ioWriterList, os.Stdout)
	}
	if setting.Server.LogFilePath != "" {
		var f, err = os.OpenFile(setting.Server.LogFilePath, os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			panic(err)
		}
		ioWriterList = append(ioWriterList, f)

		defer f.Close()
	}

	initDatabase(setting, ioWriterList)
	var router *gin.Engine = initServer(setting, ioWriterList)

	fmt.Println("Start Listen!")
	err = router.Run(fmt.Sprintf("%s:%d", setting.Server.Host, setting.Server.Port))
	if err != nil {
		panic(err)
	}

	fmt.Println("End Server")

}

func readSetting() settingoperation.SettingProperties {
	var err error
	var setting settingoperation.SettingProperties
	fmt.Println("讀取設定")
	setting, err = settingoperation.ReadConfFile()
	if err != nil {
		panic(err)
	}
	fileSaveDir = setting.Server.FileSaveDir
	avatarSaveDir = path.Join(fileSaveDir, avatarDirectory)
	imageSaveDir = path.Join(fileSaveDir, imageDirectory)

	return setting
}

func checkDictoryStruct() {
	fmt.Println("檢查目錄結構")
	var err error
	var resCode int
	var imageMDSaveDir = path.Join(imageSaveDir, imageMDDirectory)
	var pathList = []string{avatarSaveDir, imageSaveDir, imageMDSaveDir}

	for _, path := range pathList {
		resCode, err = common.CheckPath(path)
		switch resCode {
		case -1:
			panic("檢查目錄 " + path + " 時，發生 error" + err.Error())
		case 0:
			err = os.Mkdir(path, 0755)
			if err != nil {
				panic("創建 " + path + " 目錄時失敗 " + err.Error())
			}
		case 1:
			panic("創建 " + path + " 目錄時失敗，因為有相同名稱的檔案")
		case 2:
			// 已經創建
		default:
			panic("Unknow CheckPath Result Code")
		}
	}

}

func initDatabase(setting settingoperation.SettingProperties, ioWriterList []io.Writer) {
	var err error
	fmt.Println("設定資料庫")
	var logLV logger.LogLevel
	if setting.Server.DebugMode {
		gin.SetMode(gin.DebugMode)
		logLV = logger.Info
	} else {
		gin.SetMode(gin.ReleaseMode)
		logLV = logger.Error
	}

	var dsn string = setting.Database.DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.New(log.New(io.MultiWriter(ioWriterList...), "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logLV, // Log level
			IgnoreRecordNotFoundError: true,  // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,  // color
		})})
	if err != nil {
		panic("failed to connect database")
	}

	databaseoperation.SetDB(db)
	err = databaseoperation.InitDatabase()
	if err != nil {
		panic(err)
	}
}

func initServer(setting settingoperation.SettingProperties, ioWriterList []io.Writer) *gin.Engine {
	fmt.Println("設定伺服器")
	if setting.Server.Stdout {
		gin.ForceConsoleColor()
	}

	var res *gorm.DB
	var err error

	// 設定 io 輸出
	gin.DefaultWriter = io.MultiWriter(ioWriterList...)

	// 初始化 hashids
	hashIDData := hashids.NewData()
	// 從資料庫讀取 salt
	hashIDSlat, res := databaseoperation.GetSetting("HashIDSalt")
	if res.Error != nil {
		databaseoperation.CreateSetting("HashIDSalt", "")
	}
	// 沒設定過則自動生成並儲存
	if hashIDSlat == "" {
		byteArray := make([]byte, 8)
		if _, err := io.ReadFull(rand.Reader, byteArray); err != nil {
			panic(err)
		}
		hashIDSlat = base64.StdEncoding.EncodeToString(byteArray)
		res = databaseoperation.UpdateSetting("HashIDSalt", hashIDSlat)
		if res.Error != nil {
			panic(err)
		}
	}
	hashIDData.Salt = hashIDSlat
	var minLen = 5
	hashIDData.MinLength = minLen
	err = databaseoperation.SetHashIDData(hashIDSlat, minLen)
	if err != nil {
		panic(err)
	}
	hashID, err = hashids.NewWithData(hashIDData)

	// 從資料庫讀取 SessionSecretKey
	secretKeyBase64, res := databaseoperation.GetSetting("SessionSecretKey")
	if res.Error != nil {
		databaseoperation.CreateSetting("SessionSecretKey", "")
	}
	var secretKey []byte
	if secretKeyBase64 == "" {
		// 產生 64 bytes 的隨機字符
		byteArray := make([]byte, 64)
		if _, err := io.ReadFull(rand.Reader, byteArray); err != nil {
			panic(err)
		}
		secretKey = byteArray
		res = databaseoperation.UpdateSetting("SessionSecretKey", base64.StdEncoding.EncodeToString(byteArray))
		if res.Error != nil {
			panic(err)
		}
	} else {
		secretKey, err = base64.StdEncoding.DecodeString(secretKeyBase64)
		if err != nil {
			panic(err)
		}
	}

	orStr, res := databaseoperation.GetSetting("OwnerRegistered")
	if res.Error != nil {
		panic(err)
	}
	ownerRegistered, err = strconv.ParseBool(orStr)
	if err != nil {
		panic(err)
	}

	webSetting["Hostname"], res = databaseoperation.GetSetting("Hostname")
	if res.Error != nil {
		panic(err)
	}
	reaStr, res := databaseoperation.GetSetting("RequireEmailActivate")
	if res.Error != nil {
		panic(err)
	}
	webSetting["RequireEmailActivate"] = reaStr
	webSetting["SenderEmailServer"], res = databaseoperation.GetSetting("SenderEmailServer")
	if res.Error != nil {
		panic(err)
	}
	webSetting["SenderEmailAddress"], res = databaseoperation.GetSetting("SenderEmailAddress")
	if res.Error != nil {
		panic(err)
	}
	webSetting["SenderEmailUser"], res = databaseoperation.GetSetting("SenderEmailUser")
	if res.Error != nil {
		panic(err)
	}
	webSetting["SenderEmailPassword"], res = databaseoperation.GetSetting("SenderEmailPassword")
	if res.Error != nil {
		panic(err)
	}

	r := gin.Default()
	store := memstore.NewStore([]byte(secretKey))
	r.Use(sessions.Sessions("lgsc", store))
	setupRouter(r)

	return r
}

func setupRouter(r *gin.Engine) {

	var apiRouter *gin.RouterGroup = r.Group("/api")
	apiRouter.GET("/server-info", getServerInfo)
	apiRouter.GET("/dashboard/settings", getWebsiteSettings)
	apiRouter.PATCH("/dashboard/settings", editWebsiteSettings)

	apiRouter.POST("/register", register)
	apiRouter.GET("/account-activate/:token", accountActivate)
	apiRouter.POST("/signin", signin)
	apiRouter.POST("/logout", logout)
	apiRouter.GET("/me", myinfo)
	apiRouter.GET("/user/:userID", getUserInfo)
	apiRouter.GET("/user/:userID/images", getUserImages)

	apiRouter.GET("/image/:hashID", getImage)
	apiRouter.POST("/image", uploadImage)
	apiRouter.PATCH("/image/:hashID", editImage)
	apiRouter.DELETE("/image/:hashID", deleteImage)
	if !ownerRegistered {
		apiRouter.POST("/site-owner-register", siteOwnerRegister)
	}
}
