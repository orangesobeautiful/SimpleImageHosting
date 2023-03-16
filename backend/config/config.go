package config

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const AvatarDirectory = "Avatars"
const ImageDirectory = "ImagesDirect"
const ImageMDDirectory = "Medium"

var cfg CfgInfo

// CfgInfo config information
type CfgInfo struct {
	DebugMode bool

	Server struct {
		Addr           string
		FileSaveDir    string
		avatarSaveDir  string
		imageSaveDir   string
		imageMDSaveDir string
	}
	Database struct {
		DSN string
	}
	Log struct {
		Level            string
		OutputPaths      []string
		ErrorOutputPaths []string

		EnableAPIRouterLogger  bool
		EnableFileRouterLogger bool
	}
}

// Init 讀取設定並初始化
func Init(cfgPath string) (err error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(cfgPath)

	if err = v.ReadInConfig(); err != nil {
		err = fmt.Errorf("read config failed, err=%s", err)
		return
	}

	if err = v.Unmarshal(&cfg); err != nil {
		err = fmt.Errorf("unmarshal config failed, err=%s", err)
		return
	}

	if cfg.Server.FileSaveDir == "" {
		err = errors.New("FileSaveDir can not be empty")
		return
	}

	cfg.Server.avatarSaveDir = filepath.Join(cfg.Server.FileSaveDir, AvatarDirectory)
	cfg.Server.imageSaveDir = filepath.Join(cfg.Server.FileSaveDir, ImageDirectory)
	cfg.Server.imageMDSaveDir = filepath.Join(cfg.Server.imageSaveDir, ImageMDDirectory)

	return
}

// GetCfg 取得 config
func GetCfg() CfgInfo {
	return cfg
}

// AvatarSaveDir 取得頭像的儲存資料夾
func (cfg *CfgInfo) AvatarSaveDir() string {
	return cfg.Server.avatarSaveDir
}

// ImageSaveDir 取得圖片的儲存資料夾
func (cfg *CfgInfo) ImageSaveDir() string {
	return cfg.Server.imageSaveDir
}

// ImageMDSaveDir 取得中型大小的儲存資料夾
func (cfg *CfgInfo) ImageMDSaveDir() string {
	return cfg.Server.imageMDSaveDir
}

func (cfg *CfgInfo) GenZapLogger() (logger *zap.Logger, err error) {
	var zapConfig zap.Config
	if cfg.DebugMode {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	lv, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		err = fmt.Errorf("parse log level failed, err=%s", err)
		return
	}
	zapConfig.Level.SetLevel(lv)
	zapConfig.OutputPaths = cfg.Log.OutputPaths
	zapConfig.ErrorOutputPaths = cfg.Log.ErrorOutputPaths

	logger, err = zapConfig.Build()
	return
}
