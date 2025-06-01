package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Sidebar struct {
	Items         []string
	Selected      int
	Width         int
	style         lipgloss.Style
	itemStyle     lipgloss.Style
	selectedStyle lipgloss.Style
}

func NewSidebar() Sidebar {
	return Sidebar{
		Items:    []string{"# general", "# random", "# dev"},
		Selected: 0,
		Width:    20,
		style: lipgloss.NewStyle().
			Width(20).
			Border(lipgloss.RoundedBorder()).
			Padding(1, 1),
		itemStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")),
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(true),
	}
}

func (s Sidebar) Update(msg tea.Msg) (Sidebar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if s.Selected > 0 {
				s.Selected--
			}
		case "down":
			if s.Selected < len(s.Items)-1 {
				s.Selected++
			}
		}
	}
	return s, nil
}

func (s Sidebar) View() string {
	var b strings.Builder
	for i, item := range s.Items {
		if i == s.Selected {
			b.WriteString(s.selectedStyle.Render(item))
		} else {
			b.WriteString(s.itemStyle.Render(item))
		}
		b.WriteString("\n")
	}
	return s.style.Render(b.String())
}
