package ui

import (
	"fmt"
	"slacktui/components"
	"slacktui/structs"
	"slacktui/utils"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainActivityModel struct {
	chathistory    components.ChatHistory
	textarea       textarea.Model
	sidebar        *components.Sidebar
	sidebarbuttons *components.SidebarButtonView
	width          int
	height         int
	focused        string // Track which component is focused
}

type TickMsg struct{}

type UserDataLoadedMsg struct {
	Channels []structs.Channel
	DMs      []structs.DMChannel
	Err      error
}

type MessagesLoadedMsg struct {
	Messages []structs.Message
}

func loadUserDataCmd() tea.Cmd {
	return func() tea.Msg {
		channels, dms, err := utils.GetUserData()
		return UserDataLoadedMsg{Channels: channels, DMs: dms, Err: err}
	}
}

func loadMessagesCmd(channel structs.Channel) tea.Cmd {
	return func() tea.Msg {
		messages := utils.FetchChannelData(channel, 0, false)
		return MessagesLoadedMsg{messages}
	}
}

func updateDBCmd() tea.Cmd {
	return func() tea.Msg {
		success := utils.UpdateDB()
		return success
	}
}

func NewMainActivityModel() MainActivityModel {
	ta := textarea.New()
	ta.Placeholder = "Type here..."
	ta.Prompt = ""
	ta.ShowLineNumbers = false

	return MainActivityModel{
		chathistory: *components.NewChatHistory(),
		textarea:    ta,
		// With example channel and dm
		sidebar:        components.NewSidebar([]structs.Channel{}, []structs.DMChannel{}),
		sidebarbuttons: components.NewSidebarButtonView(),
	}
}

func (m MainActivityModel) Init() tea.Cmd {
	return tea.Batch(
		loadUserDataCmd(),
		updateDBCmd(),
		m.textarea.Focus(),
		tea.Tick(time.Second, func(time.Time) tea.Msg { return TickMsg{} }),
	)
}

func (m MainActivityModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var textCmd tea.Cmd
	var sidebarCmd tea.Cmd
	var chatCmd tea.Cmd

	m.sidebarbuttons, sidebarCmd = m.sidebarbuttons.Update(msg)

	switch msg := msg.(type) {
	case UserDataLoadedMsg:
		if msg.Err != nil {
			fmt.Println("Error fetching user data:", msg.Err)
			return m, nil
		}
		m.sidebar.ChannelItems = msg.Channels
		m.sidebar.DmsItems = msg.DMs
		m.sidebar.DmSelected = msg.DMs[0]
		m.sidebar.ChannelSelected = msg.Channels[0]
		m.sidebar.ReloadItems()
		// Send the msg for loading messages
		return m, loadMessagesCmd(m.sidebar.ChannelSelected)

	case MessagesLoadedMsg:
		m.chathistory.Messages = msg.Messages
		m.chathistory.ReloadMessages()
	case tea.WindowSizeMsg:
		m.width = msg.Width - 6
		m.height = msg.Height - 4

		sidebarWidth := m.width / 4
		contentWidth := m.width - sidebarWidth

		m.chathistory.Width = m.width / 4 * 3
		m.chathistory.Height = m.height - 8 // minus textarea height
		m.textarea.SetWidth(contentWidth)
		m.textarea.SetHeight(3)

		m.sidebar, sidebarCmd = m.sidebar.Update(msg)
		updatedChatHistory, _ := m.chathistory.Update(msg)
		m.chathistory = *updatedChatHistory
		cmd = m.textarea.Focus()
	case tea.KeyMsg:
		if msg.String() == "shift+tab" {
			m.sidebarbuttons.Selected += 1
			if m.sidebarbuttons.Selected >= 3 {
				m.sidebarbuttons.Selected = 0
			}
			m.sidebar.Selected += 1
			if m.sidebar.Selected >= 3 {
				m.sidebar.Selected = 0
			}
			// Run the update method to refresh the sidebar
			m.sidebar, sidebarCmd = m.sidebar.Update(msg)

		}

	case TickMsg:
		// Handle tick messages for periodic updates
		cmd = tea.Tick(time.Second*60, func(time.Time) tea.Msg {
			return TickMsg{}
		})
	}

	m.textarea, textCmd = m.textarea.Update(msg)
	return m, tea.Batch(cmd, textCmd, sidebarCmd, chatCmd)
}

func (m MainActivityModel) View() string {
	sidebarView := m.sidebar.View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		chatContentStyle.Render(m.chathistory.View()),
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
	BorderForeground(lipgloss.Color("2")).
	Padding(0, 1).
	Foreground(lipgloss.Color("111"))
