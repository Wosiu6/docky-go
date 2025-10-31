package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/wosiu6/docky-go/internal/fetcher"
)

func renderPostgres(container fetcher.ContainerInfo, width, height int) string {
	colorBorder := lipgloss.Color(colorPostgres)
	icon := "\U0001F418"
	name := baseName(container)
	name = TruncateString(name, width-4)
	var dbName, maxConn string
	if detail := container.Specific; detail != nil {
		fields := detail.DetailFields()
		dbName = fields["Database"]
		maxConn = fields["Max Conn"]
	}
	lines := []string{titleLine(icon, name, width, colorBorder)}
	if dbName != "" {
		lines = append(lines, lipgloss.NewStyle().Foreground(colorBorder).Bold(true).Render("DB: "+dbName))
	}
	lines = append(lines, statusLine(container))
	lines = append(lines, combinedStatsLine(container, "CPU: %.1f%%  MEM: %dMB"))
	if maxConn != "" {
		lines = append(lines, labelStyle.Render("Max Conn: ")+valueStyle.Render(maxConn))
	}
	lines = append(lines, imageLine(container, width), idLine(container))
	content := joinLines(lines)
	style := containerStyle.BorderForeground(colorBorder).BorderStyle(lipgloss.DoubleBorder()).Width(width)
	if height > 0 {
		style = style.Height(height)
	}
	return style.Render(content)
}
