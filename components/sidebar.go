package components

import (
	"slacktui/structs"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type KeyMap struct {
}

type Sidebar struct {
	Selected        int // Track the selected item
	viewport        viewport.Model
	width           int
	height          int
	ChannelItems    []structs.Channel
	DmsItems        []structs.DMChannel
	DmSelected      structs.DMChannel
	ChannelSelected structs.Channel
	ChannelHovered  structs.Channel
	DMHovered       structs.DMChannel
	Focused         bool // Track if the sidebar is focused
}

func NewSidebar(channelitems []structs.Channel, dmsitems []structs.DMChannel) *Sidebar {
	ourviewport := viewport.New(0, 0)
	ourviewport.KeyMap = viewport.KeyMap{}

	return &Sidebar{
		ChannelItems: channelitems,
		DmsItems:     dmsitems,
		viewport:     ourviewport,
	}
}

func (s *Sidebar) Update(msg tea.Msg) (*Sidebar, tea.Cmd) {
	var tickCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width / 4
		s.height = msg.Height - 6 // Reserve space for the footer or other UI elements
		s.viewport.Width = s.width
		s.viewport.Height = msg.Height - 8 // Adjust height to fit the viewport
	case tea.KeyMsg:
		if msg.String() == "down" && s.Focused {
			// Move down in the sidebar

			if s.Selected == 0 {
				// Handling it for channels
				for i, item := range s.ChannelItems {
					if item.ChannelId == s.ChannelHovered.ChannelId {
						if i < len(s.ChannelItems)-1 {
							s.ChannelHovered = s.ChannelItems[i+1]
						}
						break
					}
				}
			}
			if s.Selected == 1 {
				// Handling it for DMs
				for i, item := range s.DmsItems {
					if item.DmID == s.DMHovered.DmID {
						if i < len(s.DmsItems)-1 {
							s.DMHovered = s.DmsItems[i+1]
						}
						break
					}
				}
			}

			// Adjust viewport Y offset if hovered item is in the center
			if s.viewport.Height > 0 {
				center := s.viewport.Height / 2
				if s.Selected == 0 {
					for i, item := range s.ChannelItems {
						if item.ChannelId == s.ChannelHovered.ChannelId {
							if i >= center && i < len(s.ChannelItems)-center {
								s.viewport.YOffset++
							}
							break
						}
					}
				} else if s.Selected == 1 {
					for i, item := range s.DmsItems {
						if item.DmID == s.DMHovered.DmID {
							if i >= center && i < len(s.DmsItems)-center {
								s.viewport.YOffset++
							}
							break
						}
					}
				}
			}
		}
		if msg.String() == "up" && s.Focused {
			// Move up in the sidebar
			if s.Selected == 0 {
				// Handling it for channels
				for i, item := range s.ChannelItems {
					if item.ChannelId == s.ChannelHovered.ChannelId {
						if i > 0 {
							s.ChannelHovered = s.ChannelItems[i-1]
						}
						break
					}
				}
			}
			if s.Selected == 1 {
				// Handling it for DMs
				for i, item := range s.DmsItems {
					if item.DmID == s.DMHovered.DmID {
						if i > 0 {
							s.DMHovered = s.DmsItems[i-1]
						}
						break
					}
				}
			}

			// Adjust viewport Y offset if hovered item is in the center
			if s.viewport.Height > 0 {
				center := s.viewport.Height / 2
				if s.Selected == 0 {
					for i, item := range s.ChannelItems {
						if item.ChannelId == s.ChannelHovered.ChannelId {
							if i >= center && i < len(s.ChannelItems)-center {
								s.viewport.YOffset--
							}
							break
						}
					}
				} else if s.Selected == 1 {
					for i, item := range s.DmsItems {
						if item.DmID == s.DMHovered.DmID {
							if i >= center && i < len(s.DmsItems)-center {
								s.viewport.YOffset--
							}
							break
						}
					}
				}
			}

			if s.Selected == 1 {
				// Handling it for DMs
				for i, item := range s.DmsItems {
					if item.DmID == s.DMHovered.DmID {
						if i > 0 {
							s.DMHovered = s.DmsItems[i-1]
						}
						break
					}
				}
			}

			// Adjust viewport Y offset if hovered item is in the center
			if s.viewport.Height > 0 {
				center := s.viewport.Height / 2
				for i, item := range s.DmsItems {
					if item.DmID == s.DMHovered.DmID {
						if i >= center && i < len(s.DmsItems)-center {
							s.viewport.YOffset--
						}
						break
					}
				}
			}
		}
		if msg.String() == "enter" && s.Focused {
			// Make the new selected item the active one
			if s.Selected == 0 {
				s.ChannelSelected = s.ChannelHovered
			}
			if s.Selected == 1 {
				s.DmSelected = s.DMHovered
			}
			s.ReloadItems()
			s.ResetOffsetToSelected()

		}
		if msg.String() == "tab" && s.Focused {
			s.Focused = false
		}
		s.ReloadItems()

		content := s.renderItems()
		s.viewport.SetContent(content)

		var cmd tea.Cmd
		s.viewport, cmd = s.viewport.Update(msg)
		return s, tea.Batch(cmd, tickCmd)
	}
	return s, nil
}

func (s *Sidebar) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Bottom,
		sidebarStyle.Width(s.width).Height(s.height).Render(s.viewport.View()),
	)
}

func (s *Sidebar) renderItems() string {

	content := ""
	if s.Selected == 0 {
		for _, item := range s.ChannelItems {
			var aditionalStyle lipgloss.Style
			if item.ChannelId == s.ChannelHovered.ChannelId && s.Focused {
				aditionalStyle = lipgloss.NewStyle().Background(lipgloss.Color("2")).Foreground(lipgloss.Color("0"))
			}

			var icon string = "#"
			if item.IsPrivate {
				icon = ""
			}
			if item.ChannelId == s.ChannelSelected.ChannelId {

				content += lipgloss.NewStyle().
					Bold(true).Inherit(aditionalStyle).
					Foreground(lipgloss.Color("205")).
					Render("> "+icon+" "+ansi.Truncate(item.ChannelName, s.width-8, "...")) + "\n"
			} else {
				content += aditionalStyle.Render("  "+icon+" "+ansi.Truncate(item.ChannelName+" ", s.width-8, "...")) + "\n"
			}
		}
	} else if s.Selected == 1 {
		for _, item := range s.DmsItems {
			var icon string = " "
			var itemName string = ""

			if item.DmUserName != "" {
				itemName = item.DmUserName
			} else {
				itemName = item.DmUserID
			}
			var aditionalStyle lipgloss.Style
			if item.DmID == s.DMHovered.DmID && s.Focused {
				aditionalStyle = lipgloss.NewStyle().Background(lipgloss.Color("2")).Foreground(lipgloss.Color("0"))
				// White
			}

			if item.DmID == s.DmSelected.DmID {
				content += lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("205")).Inherit(aditionalStyle).
					Render("> "+icon+" "+ansi.Truncate(itemName+" ", s.width-8, "...")) + "\n"
			} else {
				content += aditionalStyle.Render("  "+icon+" "+ansi.Truncate(item.DmUserName+" ", s.width-8, "...")) + "\n"
			}
		}
	}

	return content
}

func (s *Sidebar) Init() tea.Cmd {
	return nil
}

var sidebarStyle = lipgloss.NewStyle().
	Padding(1, 2)

func (s *Sidebar) ReloadItems() {
	s.viewport.SetContent(s.renderItems())
}

func (s *Sidebar) ResetOffsetToSelected() {
	center := s.viewport.Height / 2

	if s.Selected == 0 {
		for i, item := range s.ChannelItems {
			if item.ChannelId == s.ChannelSelected.ChannelId {
				s.viewport.YOffset = i - center
				if s.viewport.YOffset < 0 {
					s.viewport.YOffset = 0
				}
				s.ChannelHovered = item
				break
			}
		}
	} else if s.Selected == 1 {
		for i, item := range s.DmsItems {
			if item.DmID == s.DmSelected.DmID {
				s.viewport.YOffset = i - center
				if s.viewport.YOffset < 0 {
					s.viewport.YOffset = 0
				}
				if i >= len(s.DmsItems)-center {
					s.viewport.YOffset = len(s.DmsItems) - s.viewport.Height
				}
				s.DMHovered = item
				break
			}
		}
	}
}
