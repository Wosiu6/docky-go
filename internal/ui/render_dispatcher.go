package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wosiu6/docky-go/internal/domain"
	"github.com/wosiu6/docky-go/internal/fetcher"
)

func (m *UiModel) renderContainer(container fetcher.ContainerInfo, width, height int) string {
	switch container.Type {
	case domain.ContainerTypePostgreSQL:
		return renderPostgres(container, width, height)
	case domain.ContainerTypeMinecraft:
		return renderMinecraft(container, width, height)
	case domain.ContainerTypeTraefik:
		return renderTraefik(container, width, height)
	case domain.ContainerTypeRedis:
		return renderRedis(container, width, height)
	case domain.ContainerTypeMinio:
		return renderMinio(container, width, height)
	case domain.ContainerTypeGrafana:
		return renderGrafana(container, width, height)
	case domain.ContainerTypePrometheus:
		return renderPrometheus(container, width, height)
	default:
		return renderGeneric(container, width, height)
	}
}

func baseName(c fetcher.ContainerInfo) string {
	name := "unnamed"
	if len(c.Names) > 0 {
		name = strings.TrimPrefix(c.Names[0], "/")
	}
	return name
}

func shortID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

func titleLine(icon, name string, width int, borderColor lipgloss.Color) string {
	return titleStyle.Background(borderColor).Width(width - 4).Render(fmt.Sprintf("%s %s", icon, name))
}

func statusLine(c fetcher.ContainerInfo) string {
	colorHex, statusIcon, statusText := StatusInfo(c.Status)
	return statusStyle.Foreground(lipgloss.Color(colorHex)).Render(fmt.Sprintf("%s %s", statusIcon, statusText))
}

func combinedStatsLine(c fetcher.ContainerInfo, format string) string {
	return statsStyle.Render(fmt.Sprintf(format, c.CPUPercent, c.Mem))
}

func imageLine(c fetcher.ContainerInfo, width int) string {
	img := TruncateString(c.Image, width-12)
	return labelStyle.Render("Image:  ") + valueStyle.Render(img)
}

func idLine(c fetcher.ContainerInfo) string {
	return labelStyle.Render("ID:     ") + valueStyle.Render(shortID(c.ID))
}

func joinLines(lines []string) string {
	var b strings.Builder
	for i, l := range lines {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(l)
	}
	return b.String()
}
