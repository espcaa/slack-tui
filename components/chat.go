package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Message struct {
	Author string
	Time   string
	Date   string
	Link   string
}

type ChatHistory struct {
	Messages []Message
	Width    int
	style    lipgloss.Style
}

func NewChatHistory() ChatHistory {
	return ChatHistory{
		Messages: []Message{{"Rowan", "10:00 AM", "2023-10-01", "https://example.com"}},
		Width:    50,
		style: lipgloss.NewStyle().
			Width(50).
			Border(lipgloss.RoundedBorder()).
			Padding(1, 1).
			Background(lipgloss.Color("235")).
			Foreground(lipgloss.Color("255")),
	}
}

func (ch ChatHistory) Update(msg tea.Msg) (ChatHistory, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ch.Width = msg.Width
		ch.style = ch.style.Width(ch.Width)
	}

	return ch, nil
}
func (ch ChatHistory) View() string {
	var b strings.Builder
	for _, m := range ch.Messages {
		b.WriteString(ch.style.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.Author+": "+m.Link,
				m.Time+" "+m.Date,
			),
		) + "\n")
	}
	return b.String()
}
