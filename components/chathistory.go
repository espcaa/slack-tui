package components

import (
	"slacktui/structs"
	"slacktui/utils"
	"strconv"

	"html"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatHistory struct {
	Messages []structs.Message
	viewport viewport.Model
	Width    int
	Height   int
	Focused  bool // Track if the chat history is focused
}

func NewChatHistory() *ChatHistory {
	return &ChatHistory{
		Messages: []structs.Message{},
		viewport: viewport.New(0, 0),
		Width:    0,
		Height:   0,
	}
}

func (c *ChatHistory) Update(msg tea.Msg) (*ChatHistory, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.Width = (msg.Width / 4 * 3) - 8
		c.Height = msg.Height - 7 // Reserve space for the footer or other UI elements
		c.viewport.Width = c.Width
		c.viewport.Height = c.Height - 2 // Adjust height to fit the viewport
		c.viewport.SetContent(c.RenderMessages())

		var cmd tea.Cmd
		c.viewport, cmd = c.viewport.Update(msg)
		return c, cmd
	case tea.KeyMsg:
		if msg.String() == "tab" && c.Focused {
			c.ReloadMessages() // Refresh chat history when Focused state changes
		}
	}
	return c, nil
}

func (c *ChatHistory) View() string {
	return c.viewport.View()
}

func (c *ChatHistory) RenderMessages() string {
	// Render the messages in a grouped format
	content := ""

	if len(c.Messages) == 0 {
		content = "No messages yet."
	} else {
		var lastSender string
		var lastTimestamp int64

		for _, msg := range c.Messages {
			var aditionalYellowText string

			if msg.ThreadBroadcast {
				aditionalYellowText = " (reply to a thread broadcasted to the channel)"
			}
			if msg.SenderName == "" {
				msg.SenderName, _ = utils.GetNameFromID(msg.SenderId, false)
			}
			// Check if the sender is the same and the time difference is within a threshold (e.g., 5 minutes)
			if msg.SenderName != lastSender || msg.Timestamp-lastTimestamp > 300 || msg.ThreadBroadcast {
				// Render the username and timestamp for a new group
				content += usernameStyle.Render(msg.SenderName) + " " + clockStyle.Render(utils.TimestampToString(strconv.FormatInt(msg.Timestamp, 10))) + " " + yellowTextStyle.Render(aditionalYellowText) + "\n"
			}
			content += messageStyle.Width(c.Width).Render(html.UnescapeString(msg.Content)) + "\n"
			lastSender = msg.SenderName
			lastTimestamp = msg.Timestamp
		}
	}
	return content
}

func (c *ChatHistory) ReloadMessages() {
	c.viewport.SetContent(c.RenderMessages())

	// Set the viewport position to the bottom
	c.viewport.GotoBottom()
}

var usernameStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
var messageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
var clockStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
var yellowTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Italic(true).Bold(true)

func (c *ChatHistory) GoToBottom() {
	// Scroll to the bottom of the chat history
	c.viewport.GotoBottom()
}

func (c *ChatHistory) AppendMessage(messages structs.Message) {
	c.Messages = append(c.Messages, messages)
	c.ReloadMessages() // Refresh the viewport to show the new message
	if !c.Focused {
		c.GoToBottom() // Automatically scroll to the bottom if not focused
	}
}

func (c *ChatHistory) DeleteMessage(messageid string) {
	// Find the message by ID and remove it from the chat history
	for i, msg := range c.Messages {
		if msg.MessageId == messageid {
			c.Messages = append(c.Messages[:i], c.Messages[i+1:]...)
			c.ReloadMessages() // Refresh the viewport to reflect the deletion
			return
		}
		if !c.Focused {
			c.GoToBottom() // Automatically scroll to the bottom if not focused
		}
	}
}

func (c *ChatHistory) ModifyMessage(messageid, content string) {
	// Find the message by ID and update its content
	for i, msg := range c.Messages {
		if msg.MessageId == messageid {
			c.Messages[i].Content = content
			c.ReloadMessages() // Refresh the viewport to reflect the modification
			return
		}
		if !c.Focused {
			c.GoToBottom() // Automatically scroll to the bottom if not focused
		}
	}
}
