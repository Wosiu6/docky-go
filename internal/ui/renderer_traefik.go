package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/wosiu6/docky-go/internal/fetcher"
)

func renderTraefik(c fetcher.ContainerInfo, w, h int) string {
	name := baseName(c)
	colorBorder := lipgloss.Color(colorTraefik)
	icon := "\U0001F6A6"
	var entrypoints string
	if d := c.Specific; d != nil {
		entrypoints = d.DetailFields()["Entrypoints"]
	}
	lines := []string{titleLine(icon, name, w, colorBorder), statusLine(c), combinedStatsLine(c, "CPU %.1f%% MEM %dMB")}
	if entrypoints != "" {
		lines = append(lines, labelStyle.Render("Entrypoints: ")+valueStyle.Render(entrypoints))
	}
	style := containerStyle.BorderForeground(colorBorder).BorderStyle(lipgloss.DoubleBorder()).Width(w)
	if h > 0 {
		style = style.Height(h)
	}
	return style.Render(joinLines(lines))
}
