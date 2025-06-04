package components

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SidebarButtonView struct {
	viewport viewport.Model
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
			{Name: "#", Keybind: "c", focused: false, Id: 0},
			{Name: "", Keybind: "d", focused: false, Id: 1},
			{Name: "", Keybind: "n", focused: false, Id: 2},
		},
		Selected: 0,
	}
}

func (s *SidebarButtonView) Update(msg tea.Msg) (*SidebarButtonView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width/4 + 2
		s.height = 8 // Reserve space for the footer or other UI elements
		s.viewport.Width = s.width
		s.viewport.Height = 8 // Adjust height to fit the viewport
	}

	content := s.renderButtons()
	s.viewport.SetContent(lipgloss.JoinHorizontal(lipgloss.Bottom, content...))

	var cmd tea.Cmd
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

func (s *SidebarButtonView) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Right,
		lipgloss.NewStyle().Width(s.width).Height(s.height).Render(s.viewport.View()),
	)
}

func (s *SidebarButtonView) Init() tea.Cmd {
	return nil
}

func (s *SidebarButtonView) renderButtons() []string {
	var buttons []string
	for _, button := range s.buttons {
		var buttonContent string = normalTextStyle.Render(button.Name + "\n" + keybindStyle.Render(button.Keybind))
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
