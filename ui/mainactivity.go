package ui

import (
	sidebar "slacktui/components"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

type MainActivityModel struct {
	viewport viewport.Model
	textarea textarea.Model
	sidebar  sidebar.Sidebar
}

func NewMainActivityModel() MainActivityModel {
	ta := textarea.New()
	ta.Placeholder = "Type here..."
	ta.SetWidth(50)            // Set initial width
	ta.SetHeight(3)            // Set initial height
	ta.ShowLineNumbers = false // Optional: Hide line numbers
	ta.Prompt = ""             // Set a custom prompt
	return MainActivityModel{
		viewport: viewport.New(80, 20),
		textarea: ta,
		sidebar:  *sidebar.NewSidebar([]string{"Item 1", "Item 2", "Item 3"}),
	}
}

func (m MainActivityModel) Init() tea.Cmd {
	return m.textarea.Focus()
}

func (m MainActivityModel) Update(msg tea.Msg) (MainActivityModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport = viewport.New(msg.Width-(msg.Width/4), msg.Height-(msg.Height/8)-2) // Adjust viewport size
		m.textarea.SetWidth(msg.Width - (msg.Width / 4))                                // Adjust textarea width
		m.textarea.SetHeight(msg.Height / 8)                                            // Keep textarea height fixed
		cmd = m.textarea.Focus()                                                        // Re-focus textarea after resize
	}

	var textCmd tea.Cmd
	m.textarea, textCmd = m.textarea.Update(msg)

	return m, tea.Batch(cmd, textCmd)
}

// Sidebar style
var sidebarStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63")).
	Background(lipgloss.Color("236")).
	Padding(1, 2).
	Width(30)

// Chat content style
var chatContentStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

// Textarea style
var textareaStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240")).
	Padding(1, 2)

func (m MainActivityModel) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebarStyle.Render(m.sidebar.View()), // Render the sidebar with updated style
		lipgloss.JoinVertical(
			lipgloss.Left,
			chatContentStyle.Render(m.viewport.View()), // Place the viewport above
			textareaStyle.Render(m.textarea.View()),    // Place the textarea at the bottom
		),
	)
}
