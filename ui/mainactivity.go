package ui

import (
	"fmt"
	"slacktui/components"
	"slacktui/structs"
	"slacktui/utils"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainActivityModel struct {
	Chathistory    *components.ChatHistory
	textarea       textarea.Model
	Sidebar        *components.Sidebar
	sidebarbuttons *components.SidebarButtonView
	width          int
	height         int
	focusedPanel   string
}

func (m *MainActivityModel) AppendMessages(newMessage structs.Message, threadbroadcast bool) {
	m.Chathistory.AppendMessage(newMessage)
}

func (m *MainActivityModel) DeleteMessage(messageID string) {
	m.Chathistory.DeleteMessage(messageID)
}

func (m *MainActivityModel) ModifyMessage(messageID, newContent string) {
	m.Chathistory.ModifyMessage(messageID, newContent)
}

func (m *MainActivityModel) GetSelectedChannelID() string {
	if m.sidebarbuttons.Selected == 0 {
		return m.Sidebar.ChannelSelected.ChannelId
	}
	if m.sidebarbuttons.Selected == 1 {
		return m.Sidebar.DmSelected.DmID
	}
	return ""
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

func loadMessagesCmd(channelid string) tea.Cmd {
	return func() tea.Msg {
		messages := utils.FetchChannelData(channelid, 0, false)
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
		Chathistory: components.NewChatHistory(),
		textarea:    ta,
		// With example channel and dm
		Sidebar:        components.NewSidebar([]structs.Channel{}, []structs.DMChannel{}),
		sidebarbuttons: components.NewSidebarButtonView(),
		focusedPanel:   "textarea",
	}
}

func initializeWebSocketCmd(m *MainActivityModel) tea.Cmd {
	return func() tea.Msg {
		err := utils.InitializeWebSocket(m)
		if err != nil {
			fmt.Println("Error initializing WebSocket:", err)
		}
		return nil
	}
}

func (m MainActivityModel) Init() tea.Cmd {
	return tea.Batch(
		loadUserDataCmd(),
		updateDBCmd(),
		m.textarea.Focus(),
		tea.Tick(time.Second, func(time.Time) tea.Msg { return TickMsg{} }),
		initializeWebSocketCmd(&m),
	)
}

func (m MainActivityModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var textCmd tea.Cmd
	var sidebarCmd tea.Cmd
	var chatCmd tea.Cmd

	m.sidebarbuttons, sidebarCmd = m.sidebarbuttons.Update(msg)
	m.Sidebar, sidebarCmd = m.Sidebar.Update(msg)
	m.Chathistory, chatCmd = m.Chathistory.Update(msg)

	switch msg := msg.(type) {
	case UserDataLoadedMsg:
		if msg.Err != nil {
			fmt.Println("Error fetching user data:", msg.Err)
			return m, nil
		}
		m.Sidebar.ChannelItems = msg.Channels
		m.Sidebar.DmsItems = msg.DMs
		m.Sidebar.DmSelected = msg.DMs[0]
		m.Sidebar.ChannelSelected = msg.Channels[0]
		m.Sidebar.ChannelHovered = msg.Channels[0]
		m.Sidebar.DMHovered = msg.DMs[0]
		m.Sidebar.ResetOffsetToSelected()
		m.Sidebar.ReloadItems()
		// Send the msg for loading messages
		return m, loadMessagesCmd(m.Sidebar.ChannelSelected.ChannelId)

	case MessagesLoadedMsg:
		m.Chathistory.Messages = msg.Messages
		m.Chathistory.ReloadMessages()
	case tea.WindowSizeMsg:

		m.width = msg.Width - 6
		m.height = msg.Height - 4

		sidebarWidth := m.width / 4
		contentWidth := m.width - sidebarWidth

		m.Chathistory.Width = m.width / 4 * 3
		m.Chathistory.Height = m.height - 8 // minus textarea height
		m.textarea.SetWidth(contentWidth)
		m.textarea.SetHeight(3)
		m.Chathistory.GoToBottom()

		cmd = m.textarea.Focus()
	case tea.KeyMsg:
		if m.focusedPanel == "chat" {
			m.Chathistory, chatCmd = m.Chathistory.Update(msg)
		}
		if msg.String() == "shift+tab" {
			m.sidebarbuttons.Selected += 1
			if m.sidebarbuttons.Selected >= 3 {
				m.sidebarbuttons.Selected = 0
			}
			m.Sidebar.Selected += 1
			if m.Sidebar.Selected >= 3 {
				m.Sidebar.Selected = 0
			}
			// Run the update method to refresh the sidebar
			m.Sidebar.ResetOffsetToSelected()
			m.Sidebar, sidebarCmd = m.Sidebar.Update(msg)
			// Load messages for the selected channel or DM
			if m.Sidebar.Selected == 0 {
				m.Sidebar.ChannelSelected = m.Sidebar.ChannelHovered
				m.Sidebar.ReloadItems()
				m.Chathistory.ReloadMessages()
				return m, loadMessagesCmd(m.Sidebar.ChannelSelected.ChannelId)
			}
			if m.Sidebar.Selected == 1 {
				m.Sidebar.DmSelected = m.Sidebar.DMHovered
				m.Sidebar.ReloadItems()
				m.Chathistory.ReloadMessages()
				return m, loadMessagesCmd(m.Sidebar.DmSelected.DmID)
			}

		}
		if msg.String() == "tab" {
			if m.focusedPanel == "chat" {
				m.Chathistory.GoToBottom()
				m.focusedPanel = "sidebar"
				m.Chathistory.Focused = false
				m.Sidebar.Focused = true
				m.Sidebar.ReloadItems()
				m.textarea.Blur()
			} else if m.focusedPanel == "sidebar" {
				m.focusedPanel = "textarea"
				m.Chathistory.Focused = false
				m.Sidebar.Focused = false
				m.textarea.Focus()
			} else if m.focusedPanel == "textarea" {
				m.focusedPanel = "chat"
				m.Chathistory.Focused = true
				m.Sidebar.Focused = false
				m.Chathistory.ReloadMessages()
				m.textarea.Blur()
			}
		}
		if msg.String() == "enter" {

			// refetch the messages for the selected channel if in sidebar or send a message

			if m.focusedPanel == "sidebar" {
				if m.sidebarbuttons.Selected == 0 {
					m.Sidebar.ChannelSelected = m.Sidebar.ChannelHovered
					m.Sidebar.ReloadItems()
					m.Chathistory.ReloadMessages()
					return m, loadMessagesCmd(m.Sidebar.ChannelSelected.ChannelId)
				} else if m.sidebarbuttons.Selected == 1 {
					m.Sidebar.DmSelected = m.Sidebar.DMHovered
					m.Sidebar.ReloadItems()
					m.Chathistory.ReloadMessages()
					return m, loadMessagesCmd(m.Sidebar.DmSelected.DmID)
				}
			} else if m.focusedPanel == "textarea" {
				message := strings.ReplaceAll(m.textarea.Value(), "\n", "")
				if message == "" {
					break
				}
				if m.Sidebar.Selected == 0 {

					utils.SendMessage(message, m.Sidebar.ChannelSelected.ChannelId)
					m.textarea.SetValue("")
				} else if m.Sidebar.Selected == 1 {
					utils.SendMessage(message, m.Sidebar.DmSelected.DmID)
					m.textarea.SetValue("")
				}
			}
		}

	case TickMsg:
		// Check if the websocket is still on and if not, reinitialize it
		if !utils.CheckWebSocketConnection() {
			return m, tea.Batch(
				initializeWebSocketCmd(&m),
				tea.Tick(time.Second*1, func(time.Time) tea.Msg {
					return TickMsg{}
				}),
			)
		}
		cmd = tea.Tick(time.Second*1, func(time.Time) tea.Msg {
			return TickMsg{}
		})
	}

	m.textarea, textCmd = m.textarea.Update(msg)
	return m, tea.Batch(cmd, textCmd, sidebarCmd, chatCmd)
}

func (m MainActivityModel) View() string {
	sidebarView := m.Sidebar.View()

	var additionalChatStyle lipgloss.Style
	if m.focusedPanel == "chat" {
		additionalChatStyle = lipgloss.NewStyle().BorderForeground(lipgloss.Color("2"))
	} else {
		additionalChatStyle = lipgloss.NewStyle()
	}
	var additionalSidebarStyle lipgloss.Style
	if m.focusedPanel == "sidebar" {
		additionalSidebarStyle = lipgloss.NewStyle().BorderForeground(lipgloss.Color("2"))
	} else {
		additionalSidebarStyle = lipgloss.NewStyle()
	}
	var additionalTextareaStyle lipgloss.Style
	if m.focusedPanel == "textarea" {
		additionalTextareaStyle = lipgloss.NewStyle().BorderForeground(lipgloss.Color("2"))
	} else {
		additionalTextareaStyle = lipgloss.NewStyle()
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		chatContentStyle.Inherit(additionalChatStyle).Render(m.Chathistory.View()),
		//textareaStyle.Width(m.textarea.Width()).Render(m.textarea.View()),
		textareaStyle.
			Inherit(additionalTextareaStyle).
			Width(m.textarea.Width()).
			Render(m.textarea.View()),
	)

	sidebarcontent := lipgloss.JoinVertical(lipgloss.Bottom,
		m.sidebarbuttons.View(),
		sidebarStyle.Inherit(additionalSidebarStyle).Render(sidebarView),
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

var sidebarStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
