package sidebar

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Sidebar struct {
	items    []string
	selected int
	viewport viewport.Model
}

func NewSidebar(items []string) *Sidebar {
	return &Sidebar{
		items:    items,
		selected: 0,
	}
}

func (s *Sidebar) SelectNext() {
	if len(s.items) == 0 {
		return
	}
	s.selected = (s.selected + 1) % len(s.items)
}

func (s *Sidebar) SelectPrevious() {
	if len(s.items) == 0 {
		return
	}
	s.selected = (s.selected - 1 + len(s.items)) % len(s.items)
}

func (s *Sidebar) GetSelected() string {
	if len(s.items) == 0 {
		return ""
	}
	return s.items[s.selected]
}

func (s *Sidebar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			s.SelectNext()
		case "k", "up":
			s.SelectPrevious()
		}
	}

	// Update the viewport if necessary
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		s.viewport.Width = msg.Width / 4
		s.viewport.Height = msg.Height
	}

	var cmd tea.Cmd
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

var sidebarStyle = lipgloss.NewStyle().
	BorderRight(true).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("2"))

func (s *Sidebar) View() string {
	var output string
	for i, item := range s.items {
		if i == s.selected {
			output += "> " + item + "\n"
		} else {
			output += "  " + item + "\n"
		}
	}
	return output + s.viewport.View()
}
func (s *Sidebar) SetSize(width, height int) {
	s.viewport.Width = width
	s.viewport.Height = height
}

func (s *Sidebar) SetViewportSize(width, height int) {
	s.viewport.Width = width
	s.viewport.Height = height
}

func (s *Sidebar) Init() tea.Cmd {
	return nil
}
