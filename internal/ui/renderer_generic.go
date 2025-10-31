package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wosiu6/docky-go/internal/fetcher"
)

func renderGeneric(container fetcher.ContainerInfo, width, height int) string {
	colorBorder := lipgloss.Color(colorGeneric)
	icon := "\U0001F4E6"
	typeLabel := string(container.Type)
	name := TruncateString(baseName(container), width-4)
	var b strings.Builder
	b.WriteString(titleLine(icon, name, width, colorBorder) + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(colorBorder).Bold(true).Render(fmt.Sprintf("\u25cf %s", typeLabel)) + "\n")
	colorHex, statusIcon, statusText := StatusInfo(container.Status)
	b.WriteString(statusStyle.Foreground(lipgloss.Color(colorHex)).Render(fmt.Sprintf("%s %s", statusIcon, statusText)) + "\n\n")
	b.WriteString(labelStyle.Render("CPU:    ") + statsStyle.Render(fmt.Sprintf("%.1f%%", container.CPUPercent)) + "\n")
	b.WriteString(labelStyle.Render("Memory: ") + statsStyle.Render(fmt.Sprintf("%d MB", container.Mem)) + "\n\n")
	image := TruncateString(container.Image, width-12)
	b.WriteString(labelStyle.Render("Image:  ") + valueStyle.Render(image) + "\n")
	b.WriteString(labelStyle.Render("ID:     ") + valueStyle.Render(shortID(container.ID)) + "\n\n")
	if detail := container.Specific; detail != nil {
		for k, v := range detail.DetailFields() {
			b.WriteString(labelStyle.Render(k+": ") + valueStyle.Render(v) + "\n")
		}
	}
	style := containerStyle.BorderForeground(colorBorder).Width(width)
	if height > 0 {
		style = style.Height(height)
	}
	return style.Render(b.String())
}
