package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

const DefaultThemeFilename = "theme.json"

type appConfig struct {
	Auth       appAuth
	Theme      appTheme
	ConfigPath string
}

type appTheme struct {
	Theme string `json:"theme,omitempty"`
}

func (mm *MainModel) loadConfig() {
	conf := appConfig{}

	path, pathErr := configPath()
	if pathErr != nil || path == "" {
		panic(fmt.Sprintf("Error loading config path: %s", pathErr.Error()))
	}
	conf.ConfigPath = path

	mm.Config = conf
}

func (mm *MainModel) loadTheme() {
	themeData, themeDataErr := loadFile(filepath.Join(mm.Config.ConfigPath, DefaultThemeFilename))
	if themeDataErr != nil {
		panic(fmt.Sprintf("Error reading theme file: %s", themeDataErr.Error()))
	}

	if len(themeData) > 0 {
		var appT appTheme
		if err := json.Unmarshal(themeData, &appT); err != nil {
			panic(fmt.Sprintf("Error unmarshalling theme json: %s", themeDataErr.Error()))
		}

		mm.Config.Theme = appT
	}

	if mm.Config.Theme.Theme != "" {
		if applyThemeErr := styles.ApplyTheme(mm.Config.Theme.Theme); applyThemeErr != nil {
			log.Printf("Failed to apply theme %s: %s", mm.Config.Theme.Theme, applyThemeErr.Error())
		}
	}
}
