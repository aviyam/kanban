package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	InitDB()
	defer CloseDB()

	config := LoadConfig()
	_ = CleanupDoneTasks(config.CleanupDoneAfterDays)
	p := tea.NewProgram(NewModel(config), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
