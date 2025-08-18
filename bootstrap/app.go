package bootstrap

import (
	"ka-cache/config"
	"ka-cache/logger"
)

const (
	CfgFilePath = "./config/config-local"
)

var (
	App *Application
)

type Application struct {
	Config *config.Config
	Logger logger.Logger
}

func init() {
	AppInit()
}

func AppInit() {
	App = &Application{}
	App.Config = loadConfig(CfgFilePath)
	App.Logger = initLogger()
}
