package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	gocache "github.com/patrickmn/go-cache"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

const DefaultThemeFilename = "theme.json"
const DefaultKeybindsFilename = "keybinds.json"
const DefaultCacheFilename = "cache.gob"

type appConfig struct {
	Auth       appAuth
	Theme      appTheme
	Keybinds   *keymaps.AppKeybinds
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
			panic(fmt.Sprintf("Error unmarshalling theme json: %s", err.Error()))
		}

		mm.Config.Theme = appT
	}

	if mm.Config.Theme.Theme != "" {
		if applyThemeErr := styles.ApplyTheme(mm.Config.Theme.Theme); applyThemeErr != nil {
			log.Printf("Failed to apply theme %s: %s", mm.Config.Theme.Theme, applyThemeErr.Error())
		}
	}
}

func (mm *MainModel) loadKeybinds() {
	keybindsData, keybindsDataErr := loadFile(filepath.Join(mm.Config.ConfigPath, DefaultKeybindsFilename))
	if keybindsDataErr != nil {
		panic(fmt.Sprintf("Error reading keybinds file: %s", keybindsDataErr.Error()))
	}

	var appK keymaps.AppKeybinds
	if len(keybindsData) > 0 {
		if err := json.Unmarshal(keybindsData, &appK); err != nil {
			panic(fmt.Sprintf("Error unmarshalling theme json: %s", err.Error()))
		}

		mm.Config.Keybinds = &appK
	} else {
		mm.Config.Keybinds = keymaps.NewDefaultAppKeymaps()

		mm.SaveKeybindsInfo()
	}
}

func (mm *MainModel) loadCache() (map[string]gocache.Item, error) {
	cacheBytes, err := loadFile(filepath.Join(mm.Config.ConfigPath, DefaultCacheFilename))
	if err != nil {
		return nil, err
	}

	cacheMap := make(map[string]gocache.Item)
	if len(cacheBytes) > 0 {
		if err := gob.NewDecoder(bytes.NewReader(cacheBytes)).Decode(&cacheMap); err != nil {
			return nil, err
		}
	}

	return cacheMap, nil
}

func (mm *MainModel) SaveKeybindsInfo() error {
	keybindsMarshal, keybindsMarshalErr := json.Marshal(mm.Config.Keybinds)
	if keybindsMarshalErr != nil {
		return keybindsMarshalErr
	}

	return saveFile(
		keybindsMarshal,
		mm.Config.ConfigPath,
		DefaultKeybindsFilename,
	)

}

func (mm *MainModel) SaveCache() error {
	file, err := os.Create(filepath.Join(mm.Config.ConfigPath, DefaultCacheFilename))
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(mm.CyberCache.Items())
}

func (mm *MainModel) SaveThemeInfo() error {
	themeMarshal, themeMarshalErr := json.Marshal(mm.Config.Theme)
	if themeMarshalErr != nil {
		return themeMarshalErr
	}

	return saveFile(
		themeMarshal,
		mm.Config.ConfigPath,
		DefaultThemeFilename,
	)
}
