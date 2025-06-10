package components

import (
	"slacktui/structs"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Sidebar struct {
	Selected        int // Track the selected item
	viewport        viewport.Model
	width           int
	height          int
	ChannelItems    []structs.Channel
	DmsItems        []structs.DMChannel
	DmSelected      structs.DMChannel
	ChannelSelected structs.Channel
}

func NewSidebar(channelitems []structs.Channel, dmsitems []structs.DMChannel) *Sidebar {
	return &Sidebar{
		ChannelItems: channelitems,
		DmsItems:     dmsitems,
	}
}

func (s *Sidebar) Update(msg tea.Msg) (*Sidebar, tea.Cmd) {
	var tickCmd tea.Cmd
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
	return s, tea.Batch(cmd, tickCmd)
}

func (s *Sidebar) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Bottom,
		sidebarStyle.Width(s.width).Height(s.height).Render(s.viewport.View()),
	)
}

func (s *Sidebar) renderItems() string {

	content := ""
	if s.Selected == 0 {
		for _, item := range s.ChannelItems {
			var icon string = "#"
			if item.ChannelId == s.ChannelSelected.ChannelId {
				content += lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("205")).
					Render("> "+icon+" "+item.ChannelName) + "\n"
			} else {
				content += "  " + icon + " " + item.ChannelName + "\n"
			}
		}
	} else if s.Selected == 1 {
		for _, item := range s.DmsItems {
			var icon string = " "
			var itemName string = ""

			if item.DmUserName != "" {
				itemName = item.DmUserName
			} else {
				itemName = item.DmUserID
			}

			if item.DmID == s.DmSelected.DmID {
				content += lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("205")).
					Render("> "+icon+" "+itemName) + "\n"
			} else {
				content += "  " + icon + " " + itemName + "\n"
			}
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

func (s *Sidebar) ReloadItems() {
	s.viewport.SetContent(s.renderItems())
}
