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
	var tickCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.Width = (msg.Width / 4 * 3) - 7
		c.Height = msg.Height - 7 // Reserve space for the footer or other UI elements
		c.viewport.Width = c.Width
		c.viewport.Height = c.Height - 2 // Adjust height to fit the viewport
	}

	content := ""
	c.viewport.SetContent(content)

	var cmd tea.Cmd
	c.viewport, cmd = c.viewport.Update(msg)
	return c, tea.Batch(cmd, tickCmd)
}

func (c *ChatHistory) View() string {
	return c.viewport.View()
}

func (c *ChatHistory) RenderMessages() string {
	// Render the messages in a simple format
	content := ""

	if len(c.Messages) == 0 {
		content = "No messages yet."
	} else {
		for _, msg := range c.Messages {
			content += usernameStyle.Render(msg.SenderName) + " -- " + utils.TimestampToString(strconv.FormatInt(msg.Timestamp, 10)) + "\n" + messageStyle.Render(html.UnescapeString(msg.Content)) + "\n"
		}
	}
	return content
}

func (c *ChatHistory) ReloadMessages() {
	c.viewport.SetContent(c.RenderMessages())
}

var usernameStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
var messageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3	"))
