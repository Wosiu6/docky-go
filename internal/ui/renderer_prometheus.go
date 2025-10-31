package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/wosiu6/docky-go/internal/fetcher"
)

func renderPrometheus(c fetcher.ContainerInfo, w, h int) string {
	name := baseName(c)
	colorBorder := lipgloss.Color(colorPrometheus)
	icon := "\U0001F525"
	var scrape string
	if d := c.Specific; d != nil {
		scrape = d.DetailFields()["Targets"]
	}
	lines := []string{titleLine(icon, name, w, colorBorder), statusLine(c), combinedStatsLine(c, "CPU %.1f%% MEM %dMB")}
	if scrape != "" {
		lines = append(lines, labelStyle.Render("Scrape Targets: ")+valueStyle.Render(scrape))
	}
	border := lipgloss.Border{Top: "\u00B7", Bottom: "\u00B7", Left: ":", Right: ":", TopLeft: "*", TopRight: "*", BottomLeft: "*", BottomRight: "*"}
	style := containerStyle.BorderForeground(colorBorder).BorderStyle(border).Width(w)
	if h > 0 {
		style = style.Height(h)
	}
	return style.Render(joinLines(lines))
}
