package models

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	settingsmodels "github.com/unremarkablegarden/cyberspace-tui-go/internal/models/settings_models"
)

type SettingsModel struct {
	activeModel tea.Model
	keybinds    keymaps.AppKeybinds
}

func NewSettingsModel(keymap keymaps.AppKeybinds) *SettingsModel {
	return &SettingsModel{
		keybinds: keymap,
	}
}

func (sm *SettingsModel) Init() tea.Cmd {
	sm.activeModel = settingsmodels.NewSettingsIndexModel(sm.keybinds)
	return nil
}

func (sm *SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Send message to active model to handle it there
	default:
		updatedModel, command := sm.activeModel.Update(msg)
		sm.activeModel = updatedModel
		return sm, command
	}

	var cmd tea.Cmd
	return sm, cmd
}

func (sm *SettingsModel) View() string {
	return sm.activeModel.View()
}
