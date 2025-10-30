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
		return renderPostgresFancy(container, width, height)
	case domain.ContainerTypeMinecraft:
		return renderMinecraftFancy(container, width, height)
	case domain.ContainerTypeTraefik:
		return renderTraefikFancy(container, width, height)
	case domain.ContainerTypeRedis:
		return renderRedisFancy(container, width, height)
	case domain.ContainerTypeMinio:
		return renderMinioFancy(container, width, height)
	case domain.ContainerTypeGrafana:
		return renderGrafanaFancy(container, width, height)
	case domain.ContainerTypePrometheus:
		return renderPrometheusFancy(container, width, height)
	default:
		return renderGenericFancy(container, width, height)
	}
}

func renderPostgresFancy(container fetcher.ContainerInfo, width, height int) string {
	borderColor := postgresColor
	icon := "\U0001F418"
	name := baseName(container)
	name = TruncateString(name, width-4)
	var dbName, maxConn string
	if detail := container.Specific; detail != nil {
		fields := detail.DetailFields()
		dbName = fields["Database"]
		maxConn = fields["Max Conn"]
	}
	var content strings.Builder
	content.WriteString(titleStyle.Background(borderColor).Width(width-4).Render(fmt.Sprintf("%s %s", icon, name)) + "\n")
	if dbName != "" {
		content.WriteString(lipgloss.NewStyle().Foreground(borderColor).Bold(true).Render("DB: "+dbName) + "\n")
	}
	colorHex, statusIcon, statusText := StatusInfo(container.Status)
	content.WriteString(statusStyle.Foreground(lipgloss.Color(colorHex)).Render(fmt.Sprintf("%s %s", statusIcon, statusText)) + "\n")
	content.WriteString(statsStyle.Render(fmt.Sprintf("CPU: %.1f%%  MEM: %dMB", container.CPUPercent, container.Mem)) + "\n")
	if maxConn != "" {
		content.WriteString(labelStyle.Render("Max Conn: ") + valueStyle.Render(maxConn) + "\n")
	}
	image := TruncateString(container.Image, width-12)
	content.WriteString(labelStyle.Render("Image:  ") + valueStyle.Render(image) + "\n")
	shortID := shortID(container.ID)
	content.WriteString(labelStyle.Render("ID:     ") + valueStyle.Render(shortID) + "\n")
	styledBox := containerStyle.BorderForeground(borderColor).BorderStyle(lipgloss.DoubleBorder()).Width(width).Height(height).Render(content.String())
	return styledBox
}

func renderMinecraftFancy(container fetcher.ContainerInfo, width, height int) string {
	borderColor := minecraftColor
	icon := "\u26CF\uFE0F"
	name := TruncateString(baseName(container), width-4)
	var players, version string
	if detail := container.Specific; detail != nil {
		fields := detail.DetailFields()
		players = fields["Players"]
		version = fields["Version"]
	}
	var content strings.Builder
	content.WriteString(titleStyle.Background(borderColor).Width(width-4).Render(fmt.Sprintf("%s %s", icon, name)) + "\n")
	if version != "" {
		content.WriteString(lipgloss.NewStyle().Foreground(borderColor).Bold(true).Render("Version: "+version) + "\n")
	}
	if players != "" {
		content.WriteString(lipgloss.NewStyle().Foreground(borderColor).Bold(true).Render("Players: "+players) + "\n")
	}
	colorHex, statusIcon, statusText := StatusInfo(container.Status)
	content.WriteString(statusStyle.Foreground(lipgloss.Color(colorHex)).Render(fmt.Sprintf("%s %s", statusIcon, statusText)) + "\n")
	content.WriteString(statsStyle.Render(fmt.Sprintf("CPU: %.1f%%  MEM: %dMB", container.CPUPercent, container.Mem)) + "\n")
	image := TruncateString(container.Image, width-12)
	content.WriteString(labelStyle.Render("Image:  ") + valueStyle.Render(image) + "\n")
	content.WriteString(labelStyle.Render("ID:     ") + valueStyle.Render(shortID(container.ID)) + "\n")
	pixelBorder := lipgloss.Border{Top: "\u2592", Bottom: "\u2592", Left: "\u2591", Right: "\u2591", TopLeft: "\u2593", TopRight: "\u2593", BottomLeft: "\u2593", BottomRight: "\u2593"}
	styledBox := containerStyle.BorderForeground(borderColor).BorderStyle(pixelBorder).Width(width).Height(height).Render(content.String())
	return styledBox
}

func renderGenericFancy(container fetcher.ContainerInfo, width, height int) string {
	borderColor := genericColor
	icon := "\U0001F4E6"
	typeLabel := string(container.Type)
	name := TruncateString(baseName(container), width-4)
	var content strings.Builder
	content.WriteString(titleStyle.Background(borderColor).Width(width-4).Render(fmt.Sprintf("%s %s", icon, name)) + "\n\n")
	typeBadge := lipgloss.NewStyle().Foreground(borderColor).Bold(true).Render(fmt.Sprintf("\u25cf %s", typeLabel))
	content.WriteString(typeBadge + "\n")
	colorHex, statusIcon, statusText := StatusInfo(container.Status)
	status := statusStyle.Foreground(lipgloss.Color(colorHex)).Render(fmt.Sprintf("%s %s", statusIcon, statusText))
	content.WriteString(status + "\n\n")
	content.WriteString(labelStyle.Render("CPU:    ") + statsStyle.Render(fmt.Sprintf("%.1f%%", container.CPUPercent)) + "\n")
	content.WriteString(labelStyle.Render("Memory: ") + statsStyle.Render(fmt.Sprintf("%d MB", container.Mem)) + "\n\n")
	image := TruncateString(container.Image, width-12)
	content.WriteString(labelStyle.Render("Image:  ") + valueStyle.Render(image) + "\n")
	content.WriteString(labelStyle.Render("ID:     ") + valueStyle.Render(shortID(container.ID)) + "\n\n")
	if detail := container.Specific; detail != nil {
		for k, v := range detail.DetailFields() {
			content.WriteString(labelStyle.Render(k+": ") + valueStyle.Render(v) + "\n")
		}
	}
	styledBox := containerStyle.BorderForeground(borderColor).Width(width).Height(height).Render(content.String())
	return styledBox
}

func renderTraefikFancy(c fetcher.ContainerInfo, w, h int) string {
	name := baseName(c)
	borderColor := lipgloss.Color("#24A1C1")
	icon := "\U0001F6A6"
	var entrypoints string
	if d := c.Specific; d != nil {
		entrypoints = d.DetailFields()["Entrypoints"]
	}
	var b strings.Builder
	b.WriteString(titleStyle.Background(borderColor).Width(w-4).Render(fmt.Sprintf("%s %s", icon, name)) + "\n")
	colorHex, statusIcon, statusText := StatusInfo(c.Status)
	b.WriteString(statusStyle.Foreground(lipgloss.Color(colorHex)).Render(fmt.Sprintf("%s %s", statusIcon, statusText)) + "\n")
	b.WriteString(statsStyle.Render(fmt.Sprintf("CPU %.1f%% MEM %dMB", c.CPUPercent, c.Mem)) + "\n")
	if entrypoints != "" {
		b.WriteString(labelStyle.Render("Entrypoints: ") + valueStyle.Render(entrypoints) + "\n")
	}
	return containerStyle.BorderForeground(borderColor).BorderStyle(lipgloss.DoubleBorder()).Width(w).Height(h).Render(b.String())
}

func renderRedisFancy(c fetcher.ContainerInfo, w, h int) string {
	name := baseName(c)
	borderColor := lipgloss.Color("#D82C20")
	icon := "\U0001F9E0"
	var mode string
	if d := c.Specific; d != nil {
		mode = d.DetailFields()["Mode"]
	}
	var b strings.Builder
	b.WriteString(titleStyle.Background(borderColor).Width(w-4).Render(fmt.Sprintf("%s %s", icon, name)) + "\n")
	colorHex, statusIcon, statusText := StatusInfo(c.Status)
	b.WriteString(statusStyle.Foreground(lipgloss.Color(colorHex)).Render(fmt.Sprintf("%s %s", statusIcon, statusText)) + "\n")
	b.WriteString(statsStyle.Render(fmt.Sprintf("CPU %.1f%% MEM %dMB", c.CPUPercent, c.Mem)) + "\n")
	if mode != "" {
		b.WriteString(labelStyle.Render("Mode: ") + valueStyle.Render(mode) + "\n")
	}
	return containerStyle.BorderForeground(borderColor).BorderStyle(lipgloss.ThickBorder()).Width(w).Height(h).Render(b.String())
}

func renderMinioFancy(c fetcher.ContainerInfo, w, h int) string {
	name := baseName(c)
	borderColor := lipgloss.Color("#FFBD2E")
	icon := "\U0001F5C4\uFE0F"
	var access, console string
	if d := c.Specific; d != nil {
		fields := d.DetailFields()
		access = fields["Access Key"]
		console = fields["Console Port"]
	}
	var b strings.Builder
	b.WriteString(titleStyle.Background(borderColor).Width(w-4).Render(fmt.Sprintf("%s %s", icon, name)) + "\n")
	colorHex, statusIcon, statusText := StatusInfo(c.Status)
	b.WriteString(statusStyle.Foreground(lipgloss.Color(colorHex)).Render(fmt.Sprintf("%s %s", statusIcon, statusText)) + "\n")
	b.WriteString(statsStyle.Render(fmt.Sprintf("CPU %.1f%% MEM %dMB", c.CPUPercent, c.Mem)) + "\n")
	if access != "" {
		b.WriteString(labelStyle.Render("Access: ") + valueStyle.Render(access) + "\n")
	}
	if console != "" {
		b.WriteString(labelStyle.Render("Console: ") + valueStyle.Render(console) + "\n")
	}
	border := lipgloss.Border{Top: "\u2550", Bottom: "\u2550", Left: "\u2551", Right: "\u2551", TopLeft: "\u2554", TopRight: "\u2557", BottomLeft: "\u255a", BottomRight: "\u255d"}
	return containerStyle.BorderForeground(borderColor).BorderStyle(border).Width(w).Height(h).Render(b.String())
}

func renderGrafanaFancy(c fetcher.ContainerInfo, w, h int) string {
	name := baseName(c)
	borderColor := lipgloss.Color("#F46800")
	icon := "\U0001F4CA"
	var plugins string
	if d := c.Specific; d != nil {
		plugins = d.DetailFields()["Plugins"]
	}
	var b strings.Builder
	b.WriteString(titleStyle.Background(borderColor).Width(w-4).Render(fmt.Sprintf("%s %s", icon, name)) + "\n")
	colorHex, statusIcon, statusText := StatusInfo(c.Status)
	b.WriteString(statusStyle.Foreground(lipgloss.Color(colorHex)).Render(fmt.Sprintf("%s %s", statusIcon, statusText)) + "\n")
	b.WriteString(statsStyle.Render(fmt.Sprintf("CPU %.1f%% MEM %dMB", c.CPUPercent, c.Mem)) + "\n")
	if plugins != "" {
		b.WriteString(labelStyle.Render("Plugins: ") + valueStyle.Render(plugins) + "\n")
	}
	return containerStyle.BorderForeground(borderColor).BorderStyle(lipgloss.RoundedBorder()).Width(w).Height(h).Render(b.String())
}

func renderPrometheusFancy(c fetcher.ContainerInfo, w, h int) string {
	name := baseName(c)
	borderColor := lipgloss.Color("#E6522C")
	icon := "\U0001F525"
	var scrape string
	if d := c.Specific; d != nil {
		scrape = d.DetailFields()["Targets"]
	}
	var b strings.Builder
	b.WriteString(titleStyle.Background(borderColor).Width(w-4).Render(fmt.Sprintf("%s %s", icon, name)) + "\n")
	colorHex, statusIcon, statusText := StatusInfo(c.Status)
	b.WriteString(statusStyle.Foreground(lipgloss.Color(colorHex)).Render(fmt.Sprintf("%s %s", statusIcon, statusText)) + "\n")
	b.WriteString(statsStyle.Render(fmt.Sprintf("CPU %.1f%% MEM %dMB", c.CPUPercent, c.Mem)) + "\n")
	if scrape != "" {
		b.WriteString(labelStyle.Render("Scrape Targets: ") + valueStyle.Render(scrape) + "\n")
	}
	border := lipgloss.Border{Top: "\u00B7", Bottom: "\u00B7", Left: ":", Right: ":", TopLeft: "*", TopRight: "*", BottomLeft: "*", BottomRight: "*"}
	return containerStyle.BorderForeground(borderColor).BorderStyle(border).Width(w).Height(h).Render(b.String())
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
