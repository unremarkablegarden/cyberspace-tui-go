package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// configPath returns the path to config folder
// tries first $CYBERSPACE_CONFIG_PATH
// defaults to ~/.cyberspace/
func configPath() (string, error) {
	dir := os.Getenv("CYBERSPACE_CONFIG_PATH")
	if dir != "" {
		return dir, nil
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		panic(fmt.Sprintf("Error getting default user home: %s", homeErr.Error()))
	}

	return filepath.Join(home, ".cyberspace"), nil
}

func loadFile(path string) ([]byte, error) {
	data, dataErr := os.ReadFile(path)
	if dataErr != nil {
		if !os.IsNotExist(dataErr) {
			return []byte{}, dataErr
		}
		return []byte{}, nil
	}

	return data, nil
}

func saveFile(data []byte, path string, filename string) error {
	if mkdirErr := os.MkdirAll(path, 0700); mkdirErr != nil {
		return mkdirErr
	}

	return os.WriteFile(filepath.Join(path, filename), data, 0600)
}

func removeConfig(path string, filename string) error {
	if rmFileErr := os.Remove(filepath.Join(path, filename)); rmFileErr != nil {
		if os.IsNotExist(rmFileErr) {
			return nil
		}
		return rmFileErr
	}

	return nil
}

func quitFilter(m tea.Model, msg tea.Msg) tea.Msg {
	if _, ok := msg.(tea.QuitMsg); !ok {
		return msg
	}

	mainModel := m.(*MainModel)

	mainModel.SaveCache()

	return msg

}
