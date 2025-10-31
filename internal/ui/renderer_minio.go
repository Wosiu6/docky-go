package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/wosiu6/docky-go/internal/fetcher"
)

func renderMinio(c fetcher.ContainerInfo, w, h int) string {
	name := baseName(c)
	colorBorder := lipgloss.Color(colorMinio)
	icon := "\U0001F5C4\uFE0F"
	var access, console string
	if d := c.Specific; d != nil {
		fields := d.DetailFields()
		access = fields["Access Key"]
		console = fields["Console Port"]
	}
	lines := []string{titleLine(icon, name, w, colorBorder), statusLine(c), combinedStatsLine(c, "CPU %.1f%% MEM %dMB")}
	if access != "" {
		lines = append(lines, labelStyle.Render("Access: ")+valueStyle.Render(access))
	}
	if console != "" {
		lines = append(lines, labelStyle.Render("Console: ")+valueStyle.Render(console))
	}
	border := lipgloss.Border{Top: "\u2550", Bottom: "\u2550", Left: "\u2551", Right: "\u2551", TopLeft: "\u2554", TopRight: "\u2557", BottomLeft: "\u255a", BottomRight: "\u255d"}
	style := containerStyle.BorderForeground(colorBorder).BorderStyle(border).Width(w)
	if h > 0 {
		style = style.Height(h)
	}
	return style.Render(joinLines(lines))
}
