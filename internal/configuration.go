package internal

import (
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/system"
)

type AppConfig struct {
	IsDarkTheme  bool `json:"darkTheme,omitempty"`
	IsLightTheme bool `json:"lightTheme,omitempty"`
}

var (
	configFile string
	configData AppConfig
)

func loadConfigurations() {
	configData = AppConfig{IsDarkTheme: false, IsLightTheme: true}
	configFile = file.JoinPath(system.HomeUserConfigDir(), "trash-app-tray-win.json")
	if file.IsFile(configFile) {
		data, err := file.ReadJsonFile[AppConfig](configFile)
		logic.ProcessError(err)
		configData = data
	} else {
		updateConfigurations()
	}
}

func updateConfigurations() {
	logger.Error(file.WriteJsonFile(configFile, configData, false))
}
