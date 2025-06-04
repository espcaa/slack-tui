package components

import (
	"slacktui/structs"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Sidebar struct {
	Items       []structs.SidebarItem
	Selected    int
	viewport    viewport.Model
	width       int
	height      int
	currentItem *structs.SidebarItem
}

func NewSidebar(items []structs.SidebarItem) *Sidebar {
	return &Sidebar{
		Items: items,
	}
}

func (s *Sidebar) Update(msg tea.Msg) (*Sidebar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width / 4
		s.height = msg.Height - 6 // Reserve space for the footer or other UI elements
		s.viewport.Width = s.width
		s.viewport.Height = msg.Height - 8 // Adjust height to fit the viewport
	}
	content := s.renderItems()
	s.viewport.SetContent(content)

	var cmd tea.Cmd
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

func (s *Sidebar) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Bottom,
		sidebarStyle.Width(s.width).Height(s.height).Render(s.viewport.View()),
	)
}

func (s *Sidebar) renderItems() string {
	if len(s.Items) == 0 {
		return "No Items available"
	}

	content := ""
	for i, item := range s.Items {
		var icon string
		if item.Type == "channel" {
			icon = "#" // Folder icon for channels
		} else if item.Type == "private_channel" {
			icon = "" // Message icon for direct messages
		} else if item.Type == "dm" {
			icon = "󰍡" // Group icon for groups

		} else {
			icon = "" // Default icon for unknown types
		}
		if i == s.Selected {
			content += lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				Render("> "+icon+" "+item.Name) + "\n"
		} else {
			content += "  " + icon + " " + item.Name + "\n"
		}
	}
	return content
}

func (s *Sidebar) Init() tea.Cmd {
	return nil
}

var sidebarStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	Padding(1, 2)
