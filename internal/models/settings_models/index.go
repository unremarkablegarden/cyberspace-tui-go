package settingsmodels

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type SettingsIndexModel struct {
	list list.Model
}

func NewSettingsIndexModel() *SettingsIndexModel {
	return &SettingsIndexModel{}
}

func (sim *SettingsIndexModel) Init() tea.Cmd {
	return nil
}

func (sim *SettingsIndexModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return sim, nil
}

func (sim *SettingsIndexModel) View() string {
	return "ESTO ES DUMMY"
}
