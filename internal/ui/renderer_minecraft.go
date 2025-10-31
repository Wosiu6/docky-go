package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/wosiu6/docky-go/internal/fetcher"
)

func renderMinecraft(container fetcher.ContainerInfo, width, height int) string {
	colorBorder := lipgloss.Color(colorMinecraft)
	icon := "\u26CF\uFE0F"
	name := TruncateString(baseName(container), width-4)
	var players, version string
	if detail := container.Specific; detail != nil {
		fields := detail.DetailFields()
		players = fields["Players"]
		version = fields["Version"]
	}
	lines := []string{titleLine(icon, name, width, colorBorder)}
	if version != "" {
		lines = append(lines, lipgloss.NewStyle().Foreground(colorBorder).Bold(true).Render("Version: "+version))
	}
	if players != "" {
		lines = append(lines, lipgloss.NewStyle().Foreground(colorBorder).Bold(true).Render("Players: "+players))
	}
	lines = append(lines, statusLine(container), combinedStatsLine(container, "CPU: %.1f%%  MEM: %dMB"), imageLine(container, width), idLine(container))
	pixelBorder := lipgloss.Border{Top: "\u2592", Bottom: "\u2592", Left: "\u2591", Right: "\u2591", TopLeft: "\u2593", TopRight: "\u2593", BottomLeft: "\u2593", BottomRight: "\u2593"}
	style := containerStyle.BorderForeground(colorBorder).BorderStyle(pixelBorder).Width(width)
	if height > 0 {
		style = style.Height(height)
	}
	return style.Render(joinLines(lines))
}
