package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/wosiu6/docky-go/internal/fetcher"
)

func renderGrafana(c fetcher.ContainerInfo, w, h int) string {
	name := baseName(c)
	colorBorder := lipgloss.Color(colorGrafana)
	icon := "\U0001F4CA"
	var plugins string
	if d := c.Specific; d != nil {
		plugins = d.DetailFields()["Plugins"]
	}
	lines := []string{titleLine(icon, name, w, colorBorder), statusLine(c), combinedStatsLine(c, "CPU %.1f%% MEM %dMB")}
	if plugins != "" {
		lines = append(lines, labelStyle.Render("Plugins: ")+valueStyle.Render(plugins))
	}
	style := containerStyle.BorderForeground(colorBorder).BorderStyle(lipgloss.RoundedBorder()).Width(w)
	if h > 0 {
		style = style.Height(h)
	}
	return style.Render(joinLines(lines))
}
