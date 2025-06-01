package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SettingsModel struct{}

func NewSettingsModel() SettingsModel {
	return SettingsModel{}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "b" { // back key handled in main model
			// nothing to do here
		}
	}
	return m, nil
}

func (m SettingsModel) View() string {
	return "Settings Screen\n\nPress 'b' to go back."
}
