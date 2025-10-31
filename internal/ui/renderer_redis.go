package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/wosiu6/docky-go/internal/fetcher"
)

func renderRedis(c fetcher.ContainerInfo, w, h int) string {
	name := baseName(c)
	colorBorder := lipgloss.Color(colorRedis)
	icon := "\U0001F9E0"
	var mode string
	if d := c.Specific; d != nil {
		mode = d.DetailFields()["Mode"]
	}
	lines := []string{titleLine(icon, name, w, colorBorder), statusLine(c), combinedStatsLine(c, "CPU %.1f%% MEM %dMB")}
	if mode != "" {
		lines = append(lines, labelStyle.Render("Mode: ")+valueStyle.Render(mode))
	}
	style := containerStyle.BorderForeground(colorBorder).BorderStyle(lipgloss.ThickBorder()).Width(w)
	if h > 0 {
		style = style.Height(h)
	}
	return style.Render(joinLines(lines))
}
