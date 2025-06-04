package ui

import (
	"slacktui/components"
	"slacktui/utils"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainActivityModel struct {
	viewport       viewport.Model
	textarea       textarea.Model
	sidebar        *components.Sidebar
	sidebarbuttons *components.SidebarButtonView
	width          int
	height         int
}

type TickMsg struct{}

func NewMainActivityModel() MainActivityModel {
	ta := textarea.New()
	ta.Placeholder = "Type here..."
	ta.Prompt = ""
	ta.ShowLineNumbers = false
	return MainActivityModel{
		viewport:       viewport.New(0, 0),
		textarea:       ta,
		sidebar:        components.NewSidebar(utils.GetChannelList()),
		sidebarbuttons: components.NewSidebarButtonView(),
	}
}

func (m MainActivityModel) Init() tea.Cmd {
	return tea.Batch(m.textarea.Focus(), tea.Tick(time.Second*60, func(time.Time) tea.Msg {
		return TickMsg{}
	}))
}

func (m MainActivityModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var textCmd tea.Cmd
	var sidebarCmd tea.Cmd

	m.sidebarbuttons, sidebarCmd = m.sidebarbuttons.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - 6
		m.height = msg.Height - 4

		sidebarWidth := m.width / 4
		contentWidth := m.width - sidebarWidth

		m.viewport.Width = contentWidth
		m.viewport.Height = m.height - 5 // minus textarea height
		m.textarea.SetWidth(contentWidth)
		m.textarea.SetHeight(3)

		m.sidebar, sidebarCmd = m.sidebar.Update(msg)
		cmd = m.textarea.Focus()

	case TickMsg:
		var user_channels = utils.GetChannelList()
		m.sidebar.Items = user_channels
		cmd = tea.Tick(time.Second*60, func(time.Time) tea.Msg {
			return TickMsg{}
		})
	}

	m.textarea, textCmd = m.textarea.Update(msg)
	return m, tea.Batch(cmd, textCmd, sidebarCmd)
}

func (m MainActivityModel) View() string {
	sidebarView := m.sidebar.View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		chatContentStyle.Width(m.viewport.Width).Render(m.viewport.View()),
		textareaStyle.Width(m.textarea.Width()).Render(m.textarea.View()),
	)

	sidebarcontent := lipgloss.JoinVertical(lipgloss.Bottom,
		m.sidebarbuttons.View(),
		sidebarView,
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebarcontent,
		content,
	)
}

// Styles
var chatContentStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	Padding(1, 2).
	Foreground(lipgloss.Color("63")).
	MarginTop(4)

var textareaStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	Padding(0, 1).
	Foreground(lipgloss.Color("111"))
