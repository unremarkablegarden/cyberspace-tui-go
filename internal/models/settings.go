package models

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	settingsmodels "github.com/unremarkablegarden/cyberspace-tui-go/internal/models/settings_models"
)

type SettingsModel struct {
	activeModel tea.Model
	keybinds    *keymaps.AppKeybinds
}

func NewSettingsModel(keymap *keymaps.AppKeybinds, setting uint8) *SettingsModel {
	sm := &SettingsModel{
		keybinds: keymap,
	}

	switch items.SettingsIndex(setting) {
	case items.SettingsIndexIndex:
		sm.activeModel = settingsmodels.NewSettingsIndexModel(keymap)
	case items.SettingsIndexKeybind:
		sm.activeModel = settingsmodels.NewSettingsKeybindsModel(keymap)
	}

	return sm
}

func (sm *SettingsModel) Init() tea.Cmd {
	return nil
}

func (sm *SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updatedModel, command := sm.activeModel.Update(msg)
	sm.activeModel = updatedModel
	return sm, command
}

func (sm *SettingsModel) View() string {
	return sm.activeModel.View()
}
