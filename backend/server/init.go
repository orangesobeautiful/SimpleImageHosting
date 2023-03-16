package server

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"sih/common"
	"sih/config"
	"sih/controller"
	"sih/models"
	"sih/models/svrsn"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var ownerRegistered bool

// Start start server
func Start() {
	err := config.Init("setting.yaml")
	if err != nil {
		log.Fatal(err)
	}
	cfg := config.GetCfg()
	logger, err := cfg.GenZapLogger()
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("check save directory")
	err = checkDirectoryStruct(&cfg)
	if err != nil {
		logger.Fatal("checkout directory failed", zap.Error(err))
	}

	logger.Info("init database")
	initDatabase(&cfg, logger)

	logger.Info("init server")
	router, err := initServer(&cfg, logger)
	if err != nil {
		logger.Fatal("initServer failed", zap.Error(err))
	}

	logger.Info("listen at " + cfg.Server.Addr)
	err = router.Run(cfg.Server.Addr)
	if err != nil {
		logger.Fatal("server run failed", zap.Error(err))
	}
}

func checkDirectoryStruct(cfg *config.CfgInfo) (err error) {
	var pathList = []string{cfg.AvatarSaveDir(), cfg.ImageSaveDir(), cfg.ImageMDSaveDir()}

	for _, p := range pathList {
		var resCode int
		resCode, err = common.CheckPath(p)
		switch resCode {
		case -1:
			return
		case 0:
			const dirPerm = 0755
			err = os.MkdirAll(p, dirPerm)
			if err != nil {
				err = fmt.Errorf("mkdir %s failed, err=%s", p, err)
				return
			}
		case 1:
			err = fmt.Errorf("%s is not director", p)
		case 2:
			// 已經創建
		default:
			panic("Unknow CheckPath Result Code")
		}
	}

	return
}

func initDatabase(cfg *config.CfgInfo, logger *zap.Logger) {
	var err error
	var gormCfg = &gorm.Config{}

	if cfg.DebugMode {
		gormCfg.Logger = gormLogger.New(log.New(
			io.MultiWriter(os.Stderr), "\r\n", log.LstdFlags),
			gormLogger.Config{
				LogLevel:                  gormLogger.Info, // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,            // color
			})
	}

	var dsn = cfg.Database.DSN
	db, err := gorm.Open(mysql.Open(dsn), gormCfg)
	if err != nil {
		panic("failed to connect database")
	}

	err = models.Init(db, logger)
	if err != nil {
		panic(err)
	}
}

func initServer(cfg *config.CfgInfo, logger *zap.Logger) (engine *gin.Engine, err error) {
	// 從資料庫讀取 SessionSecretKey
	secretKeyBase64 := models.SvrSettingGet(svrsn.SessSecretKey)

	var secretKey []byte
	if secretKeyBase64 == "" {
		// 產生 64 bytes 的隨機字符
		const secretKeyLen = 64
		byteArray := make([]byte, secretKeyLen)
		if _, err = io.ReadFull(rand.Reader, byteArray); err != nil {
			return
		}
		secretKey = byteArray
		err = models.SvrSettingUpdate(svrsn.SessSecretKey, base64.StdEncoding.EncodeToString(byteArray))
		if err != nil {
			return
		}
	} else {
		secretKey, err = base64.StdEncoding.DecodeString(secretKeyBase64)
		if err != nil {
			return
		}
	}

	orStr := models.SvrSettingGet(svrsn.OwnerRegistered)
	ownerRegistered, err = strconv.ParseBool(orStr)
	if err != nil {
		return
	}

	r := gin.New()
	if cfg.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(gin.Recovery())
	store := cookie.NewStore(secretKey)
	r.Use(sessions.Sessions("sihsi", store))
	setupRouter(r, cfg, logger)

	engine = r
	return engine, nil
}

func setupRouter(r *gin.Engine, cfg *config.CfgInfo, logger *zap.Logger) {
	// api router

	var apiRouter = r.Group("/api")
	if cfg.Log.EnableAPIRouterLogger {
		apiRouter.Use(ginzap.Ginzap(logger, "", false))
	}
	apiRouter.GET("/server-info", controller.GetServerInfo)
	apiRouter.GET("/dashboard/settings", controller.GetWebsiteSettings)
	apiRouter.PATCH("/dashboard/settings", controller.EditWebsiteSettings)

	apiRouter.POST("/register", controller.Register)
	apiRouter.GET("/account-activate/:token", controller.AccountActivate)
	apiRouter.POST("/signin", controller.Signin)
	apiRouter.POST("/logout", controller.Logout)
	apiRouter.GET("/me", controller.Myinfo)
	apiRouter.GET("/user/:userID", controller.GetUserInfo)
	apiRouter.GET("/user/:userID/images", controller.GetUserImages)

	apiRouter.GET("/image/:hashID", controller.GetImage)
	apiRouter.POST("/image", controller.UploadImage)
	apiRouter.PATCH("/image/:hashID", controller.EditImage)
	apiRouter.DELETE("/image/:hashID", controller.DeleteImage)
	if !ownerRegistered {
		apiRouter.POST("/site-owner-register", controller.SiteOwnerRegister)
	}

	// file router

	fP, _ := filepath.Abs(cfg.Server.FileSaveDir)
	rootGroup := r.Group("/")
	if cfg.Log.EnableFileRouterLogger {
		rootGroup.Use(ginzap.Ginzap(logger, "", false))
	}

	rootGroup.Static("/ImagesDirect/", filepath.Join(fP, "ImagesDirect"))
	rootGroup.Static("/Avatars/", filepath.Join(fP, "Avatars"))
}
