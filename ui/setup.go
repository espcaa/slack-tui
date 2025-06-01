package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"slacktui/config"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SetupModel struct {
	textInput textinput.Model
	progress  progress.Model
	done      bool
	loading   bool
	errorText string
}

func NewSetupModel() SetupModel {
	ti := textinput.New()
	ti.Placeholder = "Enter something"
	ti.Focus()
	ti.CharLimit = 156
	// Set max width for the text input
	ti.Width = 30

	p := progress.New(progress.WithSolidFill("3"))

	return SetupModel{
		textInput: ti,
		progress:  p,
	}
}

func (m SetupModel) Init() tea.Cmd {
	return textinput.Blink
}

type AuthTestResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func checkSlackToken(token string) tea.Msg {
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://slack.com/api/auth.test", nil)
	if err != nil {
		return errMsg{err}
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return errMsg{err}
	}
	defer resp.Body.Close()

	var result AuthTestResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return errMsg{err}
	}

	if result.Ok {
		return successMsg{}
	}
	return errMsg{fmt.Errorf(result.Error)}
}

type errMsg struct{ err error }
type successMsg struct{}

type tickMsg time.Time

func (m SetupModel) Update(msg tea.Msg) (SetupModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textInput.Width = msg.Width - 10

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			m.progress = progress.New(progress.WithSolidFill("3")) // reset progress
			m.progress.SetPercent(0.0)
			m.loading = true
			token := m.textInput.Value()
			if token == "" {
				m.errorText = "Token cannot be empty."
				m.loading = false
				return m, nil
			}
			m.errorText = ""
			return m, tea.Batch(tickCmd(), func() tea.Msg {
				return checkSlackToken(token)
			})

		}

	case successMsg:
		// Save the token to the file

		token := m.textInput.Value()
		var cfg, err = config.LoadConfig()
		if err != nil {
			return m, func() tea.Msg {
				return errMsg{fmt.Errorf("failed to load config: %w", err)}
			}
		}
		cfg.SlackToken = token
		if err := config.SaveConfig(cfg); err != nil {
			return m, func() tea.Msg {
				return errMsg{fmt.Errorf("failed to save config: %w", err)}
			}
		} else {
			m.done = true
		}

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			m.loading = false
		}
		if m.loading {
			cmd := m.progress.IncrPercent(1.0)
			return m, tea.Batch(tickCmd(), cmd)
		}

	case errMsg:
		m.done = false
		m.errorText = msg.err.Error()
	}

	m.textInput, cmd = m.textInput.Update(msg)
	if m.loading {
		var updated tea.Model
		updated, cmd = m.progress.Update(msg)
		m.progress = updated.(progress.Model)
	}
	return m, cmd
}

var style = lipgloss.NewStyle().
	MarginBottom(1).
	MarginLeft(1).
	Foreground(lipgloss.Color("11")).
	Border(lipgloss.NormalBorder())

var fancystyle = lipgloss.NewStyle().
	PaddingTop(2).
	PaddingBottom(1).
	Italic(true).
	Bold(true).
	Foreground(lipgloss.Color("10"))

var subtitlestyle = lipgloss.NewStyle().
	Italic(true).
	Foreground(lipgloss.Color("2"))

var errorStyle = lipgloss.NewStyle().
	Italic(true).
	Foreground(lipgloss.Color("1"))

var keyStyle = lipgloss.NewStyle().
	Bold(true)

func (m SetupModel) View() string {
	if m.loading {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			subtitlestyle.PaddingTop(2).PaddingBottom(2).Render("Checking token..."),
			m.progress.View(),
		)
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		fancystyle.Render("Welcome to Slack TUI!"),
		subtitlestyle.Render("Please enter your Slack user token to get started:"),
		style.Render(m.textInput.View()),
		fmt.Sprintf("Press %s to continue.", keyStyle.Render("Enter")),
		errorStyle.Render(m.errorText),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
