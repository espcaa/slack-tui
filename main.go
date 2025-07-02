package main

import (
	"fmt"
	"log"

	"github.com/espcaa/slack-tui/config"
	"github.com/espcaa/slack-tui/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var startView string = "setup"

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	} else if cfg.SlackToken != "" {
		startView = "main"
	}

	p := tea.NewProgram(ui.NewMainModel(startView))
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
