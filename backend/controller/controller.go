package controller

import (
	"sih/config"

	"go.uber.org/zap"
)

const (
	userkey = "user"
)

var logger *zap.Logger
var cfg config.CfgInfo

// New return controller
func Init(logger *zap.Logger) {
	logger = logger
	cfg = config.GetCfg()
}
