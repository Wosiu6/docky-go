package ui

import "github.com/charmbracelet/lipgloss"

var (
	runningColor    = lipgloss.Color("#00FF00")
	stoppedColor    = lipgloss.Color("#FF0000")
	pausedColor     = lipgloss.Color("#FFA500")
	restartingColor = lipgloss.Color("#FFFF00")
	createdColor    = lipgloss.Color("#00BFFF")

	postgresColor  = lipgloss.Color("#336791")
	minecraftColor = lipgloss.Color("#62B47A")
	portainerColor = lipgloss.Color("#13BEF9")
	genericColor   = lipgloss.Color("#874BFD")

	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			MarginRight(1).
			MarginBottom(1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(0, 1).
			MarginBottom(1)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	statsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF0000")).
			Padding(1, 2)

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)
)
