package controller

import (
	"sih/config"
)

const (
	userkey = "user"
)

var cfg config.CfgInfo

// New return controller
func Init() {
	cfg = config.GetCfg()
}
