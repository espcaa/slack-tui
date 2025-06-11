package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SidebarButtonView struct {
	width    int
	height   int
	buttons  []SidebarButton
	Selected int
}

type SidebarButton struct {
	Name    string
	Keybind string
	focused bool
	Id      int
}

func NewSidebarButtonView() *SidebarButtonView {

	return &SidebarButtonView{
		width:  20, // Default width
		height: 8,  // Default height
		buttons: []SidebarButton{
			{Name: "#", Keybind: "channels", focused: false, Id: 0},
			{Name: "", Keybind: "dms", focused: false, Id: 1},
			{Name: "", Keybind: "notifs", focused: false, Id: 2},
		},
		Selected: 0,
	}
}

func (s *SidebarButtonView) Update(msg tea.Msg) (*SidebarButtonView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width/4 + 2
		s.height = 8 // Reserve space for the footer or other UI elements
	}

	return s, nil
}

func (s *SidebarButtonView) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Right, s.renderButtons()...)
}

func (s *SidebarButtonView) Init() tea.Cmd {
	return nil
}

func (s *SidebarButtonView) renderButtons() []string {
	var buttons []string
	for _, button := range s.buttons {
		var buttonContent string = lipgloss.JoinVertical(
			lipgloss.Center,
			normalTextStyle.Render(button.Name),
			keybindStyle.Render(button.Keybind),
		)
		if s.Selected == button.Id {
			buttons = append(buttons, buttonStyle.Width((s.width-4)/3).Foreground(lipgloss.Color("205")).Render(buttonContent))
		} else {
			buttons = append(buttons, buttonStyle.Width((s.width-4)/3).Render(buttonContent))
		}
	}

	return buttons
}

var buttonStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	MarginTop(4).
	Align(lipgloss.Center, lipgloss.Center)

var keybindStyle = lipgloss.NewStyle().
	Italic(true)

var normalTextStyle = lipgloss.NewStyle().
	Bold(true)
