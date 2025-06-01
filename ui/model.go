package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	currentView string
	setup       SetupModel
	mainAct     MainActivityModel
	settings    SettingsModel
}

func NewMainModel(startView string) MainModel {
	return MainModel{
		currentView: startView,
		setup:       NewSetupModel(),
		mainAct:     NewMainActivityModel(),
		settings:    NewSettingsModel(),
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	switch m.currentView {
	case "setup":
		setup, cmd := m.setup.Update(msg)
		m.setup = setup
		if setup.done && !setup.loading {
			m.currentView = "main"
		}
		return m, cmd

	case "main":
		mainAct, cmd := m.mainAct.Update(msg)
		m.mainAct = mainAct
		// example: switch to settings on "s"

		return m, cmd

	case "settings":
		settings, cmd := m.settings.Update(msg)
		m.settings = settings
		// example: back to main on "b"

		return m, cmd
	}
	return m, nil
}

func (m MainModel) View() string {
	switch m.currentView {
	case "setup":
		return m.setup.View()
	case "main":
		return m.mainAct.View()
	case "settings":
		return m.settings.View()
	}
	return fmt.Sprintf("Unknown view: %s", m.currentView)
}
