package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var dockyLogo = []string{
	"    .___             __                                   ",
	"  __| _/____   ____ |  | _____.__.           ____   ____  ",
	" / __ |/  _ \\_/ ___\\|  |/ <   |  |  ______  / ___\\ /  _ \\ ",
	"/ /_/ (  <_> )  \\___|    < \\___  | /_____/ / /_/  >  <_> )",
	"\\____ |\\____/ \\___  >__|_ \\ ____|         \\___  / \\____/ ",
	"     \\/           \\/     \\ \\/             /_____/          ",
}

func logo() string {
	var lines []string
	lines = append(lines, dockyLogo...)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FDF500")).
		Bold(true).
		Render(strings.Join(lines, "\n"))
}
