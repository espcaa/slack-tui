package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"slacktui/config"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SetupModel struct {
	textInput    textinput.Model
	progress     progress.Model
	done         bool
	loading      bool
	errorText    string
	width        int
	checkinginfo bool
	proceeding   bool // Indicates if the user has confirmed the information
}

var tempData struct {
	Userworkspace string
	Usertoken     string
	Userusername  string
	UserteamID    string
}

func NewSetupModel() SetupModel {
	ti := textinput.New()
	ti.Placeholder = "Your magic string here!"
	// Set max width for the text input
	ti.Width = 30
	ti.Focus()

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
	Ok                  bool   `json:"ok"`
	URL                 string `json:"url"`
	Team                string `json:"team"`
	User                string `json:"user"`
	TeamID              string `json:"team_id"`
	UserID              string `json:"user_id"`
	IsEnterpriseInstall bool   `json:"is_enterprise_install"`
	Error               string `json:"error,omitempty"`
}

func checkSlackToken(input string) tea.Msg {
	// Split the input string into cookies and token
	cookies := ""
	token := ""

	parts := strings.Split(input, "||")
	if len(parts) == 2 {
		cookies = parts[0]
		token = parts[1]
	} else if len(parts) == 1 {
		return errMsg{fmt.Errorf("Magic string is invalid")}
	}

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://slack.com/api/auth.test", nil)
	if err != nil {
		return errMsg{err}
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Cookie", cookies)
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
		// If the response is ok, we can save the token, workspace, and username
		tempData.Usertoken = parts[1]
		tempData.Userworkspace = result.URL
		tempData.Userusername = result.User
		tempData.UserteamID = result.TeamID

		return successMsg{}
	}
	return errMsg{fmt.Errorf("A network error occured. Either your magic string is invalid or your internet connection is poor/unavailable.\nError: %s", result.Error)}
}

type errMsg struct{ err error }
type successMsg struct{}

type tickMsg time.Time

func (m SetupModel) Update(msg tea.Msg) (SetupModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textInput.Width = msg.Width - 10
		m.progress.Width = msg.Width - 10
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.checkinginfo {
				m.checkinginfo = false
				m.textInput.SetValue("")
				m.errorText = ""
				return m, nil
			} else {
				return m, tea.Quit
			}
		case "enter":
			if !m.checkinginfo {
				m.progress = progress.New(progress.WithSolidFill("3")) // reset progress
				m.progress.SetPercent(0.0)
				m.loading = true
				if m.textInput.Value() == "" {
					m.errorText = "Magic string cannot be empty."
					m.loading = false
					return m, nil
				}
				m.errorText = ""
				return m, tea.Batch(tickCmd(), func() tea.Msg {
					return checkSlackToken(m.textInput.Value())
				})
			} else if m.checkinginfo && !m.loading {
				m.proceeding = true
				m.loading = true
				m.progress = progress.New(progress.WithSolidFill("3")) // reset progress
				m.progress.SetPercent(0.0)
				// Handle saving of the data in config file
				var cfg, err = config.LoadConfig()
				if err != nil {
					return m, func() tea.Msg {
						return errMsg{fmt.Errorf("failed to load config: %w", err)}
					}
				}
				cfg.SlackToken = tempData.Usertoken
				cfg.WorkspaceID = tempData.UserteamID
				cfg.WorkspaceURL = tempData.Userworkspace
				cfg.Cookies = strings.Split(m.textInput.Value(), "||")[0]

				if err := config.SaveConfig(cfg); err != nil {
					return m, func() tea.Msg {
						return errMsg{fmt.Errorf("failed to save config: %w", err)}
					}
				}
				return m, tea.Batch(tickCmd())

			}

		}

	case successMsg:
		m.checkinginfo = true

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			m.loading = false
			if m.checkinginfo && m.proceeding {
				m.done = true
			}
		}
		if m.loading {
			cmd := m.progress.IncrPercent(1.0)
			return m, tea.Batch(tickCmd(), cmd)
		}

	case errMsg:
		m.done = false
		m.errorText = "A network error occured (This is either due to no or a bad internet connection, or an invalid magic link)."
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
	Bold(true).
	Foreground(lipgloss.Color("10"))

var errorStyle = lipgloss.NewStyle().
	Italic(true).
	Foreground(lipgloss.Color("1"))

var keyStyle = lipgloss.NewStyle().
	Bold(true)

func (m SetupModel) View() string {
	var subtitlestyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("5")).
		Width(m.width - 10)

	if m.checkinginfo && !m.loading {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			fancystyle.Render(`
 /\_/\
( o.o ) < Seems like it worked!!!
 > ^ <
 `),
			subtitlestyle.Render(fmt.Sprintf("Workspace: %s", tempData.Userworkspace)),
			subtitlestyle.Render(fmt.Sprintf("Username: %s", tempData.Userusername)),
			subtitlestyle.Render(fmt.Sprintf("Team ID: %s", tempData.UserteamID)),
			subtitlestyle.MarginBottom(2).Render(fmt.Sprintf("Token: %s", tempData.Usertoken)),

			fmt.Sprintf("Press %s to continue or %s to cancel.", keyStyle.Render("Enter"), keyStyle.Render("Esc")),
		)
	}

	if m.loading {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			subtitlestyle.PaddingTop(2).PaddingBottom(2).Render("Doing some work..."),
			m.progress.View(),
		)
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		fancystyle.Render(`
 /\_/\
( o.o ) < Welcome to SlackTUI!
 > ^ <
 `),
		subtitlestyle.Render("To get started we need quite a few things from slack!"),
		subtitlestyle.Render("To get all of that, just install our chrome extension and copy paste the magic string it gives you here:"),
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
