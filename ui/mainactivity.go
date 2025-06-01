package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type MainActivityModel struct{}

func NewMainActivityModel() MainActivityModel {
	return MainActivityModel{}
}

func (m MainActivityModel) Init() tea.Cmd {
	return nil
}

func (m MainActivityModel) Update(msg tea.Msg) (MainActivityModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "s" { // switch to settings handled in main model
			// nothing here
		}
	}
	return m, nil
}

func (m MainActivityModel) View() string {
	return "Main Activity Screen\n\nPress 's' to open settings."
}
