package ui

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wosiu6/docky-go/fetcher"
)

type UiModel struct {
	fetcher  *fetcher.Fetcher
	lastErr  error
	items    []fetcher.ContainerInfo
	loading  bool
	termSize tea.WindowSizeMsg
}

func New(fetcher *fetcher.Fetcher) *UiModel { return &UiModel{fetcher: fetcher} }

func (m *UiModel) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, tickCmd())
}

func tickCmd() tea.Cmd { return tea.Tick(time.Second*2, func(t time.Time) tea.Msg { return t }) }

func (m *UiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.termSize = msg
		return m, nil

	case time.Time:
		m.loading = true
		var cmds []tea.Cmd
		cmds = append(cmds, tea.ClearScreen)

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		items, err := m.fetcher.FetchAll(ctx)

		if err != nil {
			m.lastErr = err
			m.items = nil
		} else {
			sort.Slice(items, func(i, j int) bool {
				nameI, nameJ := "unnamed", "unnamed"
				if len(items[i].Names) > 0 {
					nameI = items[i].Names[0]
				}
				if len(items[j].Names) > 0 {
					nameJ = items[j].Names[0]
				}
				return nameI < nameJ
			})
			m.items = items
			m.lastErr = nil
		}

		m.loading = false
		cmds = append(cmds, tickCmd())
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

var (
	runningColor    = lipgloss.Color("#00FF00")
	stoppedColor    = lipgloss.Color("#FF0000")
	pausedColor     = lipgloss.Color("#FFA500")
	restartingColor = lipgloss.Color("#FFFF00")
	createdColor    = lipgloss.Color("#00BFFF")

	postgresColor  = lipgloss.Color("#336791")
	minecraftColor = lipgloss.Color("#62B47A")
	portainerColor = lipgloss.Color("#13BEF9")
	genericColor   = lipgloss.Color("#874BFD")

	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			MarginRight(1).
			MarginBottom(1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(0, 1).
			MarginBottom(1)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	statsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF0000")).
			Padding(1, 2)

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)
)

func (m *UiModel) View() string {
	if m.lastErr != nil {
		return errorStyle.Render(fmt.Sprintf("❌ Error: %v", m.lastErr))
	}

	if len(m.items) == 0 {
		renderString := centeredLogo() + "\n\n🐳 No containers found. Waiting for containers..."
		return emptyStyle.Render(renderString)
	}

	width := 120
	height := 30

	numContainers := len(m.items)
	var cols, rows int

	switch {
	case numContainers == 1:
		cols, rows = 1, 1
	case numContainers == 2:
		cols, rows = 2, 1
	case numContainers <= 4:
		cols, rows = 2, 2
	case numContainers <= 6:
		cols, rows = 3, 2
	default:
		cols, rows = 3, 3
	}

	boxWidth := (width / cols) - 4
	boxHeight := (height / rows) - 4

	var grid []string
	var currentRow []string

	for i, container := range m.items {
		if i >= cols*rows {
			break
		}

		box := m.renderContainer(container, boxWidth, boxHeight)
		currentRow = append(currentRow, box)

		if len(currentRow) == cols || i == len(m.items)-1 {
			grid = append(grid, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
			currentRow = []string{}
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, grid...)
}

func (m *UiModel) renderContainer(container fetcher.ContainerInfo, width, height int) string {
	name := "unnamed"
	if len(container.Names) > 0 {
		name = strings.TrimPrefix(container.Names[0], "/")
	}

	if len(name) > width-4 {
		name = name[:width-7] + "..."
	}

	borderColor := genericColor
	icon := "📦"
	typeLabel := string(container.Type)

	switch container.Type {
	case fetcher.TypePostgreSQL:
		borderColor = postgresColor
		icon = "🐘"
		typeLabel = "PostgreSQL"
	case fetcher.TypeMinecraft:
		borderColor = minecraftColor
		icon = "⛏️"
		typeLabel = "Minecraft"
	case fetcher.TypePortainer:
		borderColor = portainerColor
		icon = "🐋"
		typeLabel = "Portainer"
	case fetcher.TypeGeneric:
		icon = "📦"
		typeLabel = "Container"
	}

	var content strings.Builder

	title := titleStyle.
		Background(borderColor).
		Width(width - 4).
		Render(fmt.Sprintf("%s %s", icon, name))
	content.WriteString(title + "\n\n")

	typeBadge := lipgloss.NewStyle().
		Foreground(borderColor).
		Bold(true).
		Render(fmt.Sprintf("● %s", typeLabel))
	content.WriteString(typeBadge + "\n")

	statusColor := stoppedColor
	statusIcon := "⭘"
	statusText := strings.ToUpper(container.Status)

	switch strings.ToLower(container.Status) {
	case "running":
		statusColor = runningColor
		statusIcon = "●"
	case "paused":
		statusColor = pausedColor
		statusIcon = "❚❚"
	case "restarting":
		statusColor = restartingColor
		statusIcon = "↻"
	case "exited":
		statusColor = stoppedColor
		statusIcon = "■"
	case "created":
		statusColor = createdColor
		statusIcon = "○"
	case "dead":
		statusColor = stoppedColor
		statusIcon = "✗"
	}

	status := statusStyle.
		Foreground(statusColor).
		Render(fmt.Sprintf("%s %s", statusIcon, statusText))
	content.WriteString(status + "\n\n")

	content.WriteString(labelStyle.Render("CPU:    ") +
		statsStyle.Render(fmt.Sprintf("%.1f%%", container.CPUPercent)) + "\n")
	content.WriteString(labelStyle.Render("Memory: ") +
		statsStyle.Render(fmt.Sprintf("%d MB", container.Mem)) + "\n\n")

	image := container.Image
	if len(image) > width-12 {
		image = image[:width-15] + "..."
	}
	content.WriteString(labelStyle.Render("Image:  ") + valueStyle.Render(image) + "\n")

	shortID := container.ID
	if len(shortID) > 12 {
		shortID = shortID[:12]
	}
	content.WriteString(labelStyle.Render("ID:     ") + valueStyle.Render(shortID) + "\n")

	content.WriteString("\n")
	switch container.Type {
	case fetcher.TypePostgreSQL:
		if container.PostgreSql != nil {
			content.WriteString(m.renderPostgreSqlInfo(container.PostgreSql, width))
		}
	case fetcher.TypeMinecraft:
		if container.Minecraft != nil {
			content.WriteString(m.renderMinecraftInfo(container.Minecraft, width))
		}
	case fetcher.TypePortainer:
		if container.Portainer != nil {
			content.WriteString(m.renderPortainerInfo(container.Portainer, width))
		}
	}

	styledBox := containerStyle.
		BorderForeground(borderColor).
		Width(width).
		Height(height).
		Render(content.String())

	return styledBox
}

func (m *UiModel) renderPostgreSqlInfo(pg *fetcher.PostgreSqlContainerInfo, width int) string {
	var s strings.Builder

	s.WriteString(lipgloss.NewStyle().Foreground(postgresColor).Bold(true).Render("PostgreSQL Info") + "\n")

	if pg.Port > 0 {
		s.WriteString(labelStyle.Render("Port:     ") + valueStyle.Render(fmt.Sprintf("%d", pg.Port)) + "\n")
	}
	if pg.Database != "" {
		s.WriteString(labelStyle.Render("Database: ") + valueStyle.Render(pg.Database) + "\n")
	}
	if pg.User != "" {
		s.WriteString(labelStyle.Render("User:     ") + valueStyle.Render(pg.User) + "\n")
	}
	if pg.SSLMode != "" {
		s.WriteString(labelStyle.Render("SSL Mode: ") + valueStyle.Render(pg.SSLMode) + "\n")
	}
	if pg.MaxConnections > 0 {
		s.WriteString(labelStyle.Render("Max Conn: ") + valueStyle.Render(fmt.Sprintf("%d", pg.MaxConnections)) + "\n")
	}
	if pg.PGData != "" {
		s.WriteString(labelStyle.Render("Volume:     ") + valueStyle.Render(pg.PGData) + "\n")
	}

	return s.String()
}

func (m *UiModel) renderMinecraftInfo(mc *fetcher.MinecraftContainerInfo, width int) string {
	var s strings.Builder

	s.WriteString(lipgloss.NewStyle().Foreground(minecraftColor).Bold(true).Render("Minecraft Info") + "\n")

	if mc.Port > 0 {
		s.WriteString(labelStyle.Render("Port:     ") + valueStyle.Render(fmt.Sprintf("%d", mc.Port)) + "\n")
	}
	if mc.Version != "" {
		s.WriteString(labelStyle.Render("Version:  ") + valueStyle.Render(mc.Version) + "\n")
	}
	if mc.ServerType != "" {
		s.WriteString(labelStyle.Render("Type:     ") + valueStyle.Render(mc.ServerType) + "\n")
	}
	if mc.Difficulty != "" {
		s.WriteString(labelStyle.Render("Difficulty:") + valueStyle.Render(mc.Difficulty) + "\n")
	}
	if mc.MaxPlayers > 0 {
		s.WriteString(labelStyle.Render("Players:  ") +
			valueStyle.Render(fmt.Sprintf("%d/%d", mc.OnlinePlayers, mc.MaxPlayers)) + "\n")
	}

	return s.String()
}

func (m *UiModel) renderPortainerInfo(pt *fetcher.PortainerContainerInfo, width int) string {
	var s strings.Builder

	s.WriteString(lipgloss.NewStyle().Foreground(portainerColor).Bold(true).Render("Portainer Info") + "\n")

	if pt.Port > 0 {
		s.WriteString(labelStyle.Render("Port:    ") + valueStyle.Render(fmt.Sprintf("%d", pt.Port)) + "\n")
	}
	if pt.Edition != "" {
		s.WriteString(labelStyle.Render("Edition: ") + valueStyle.Render(pt.Edition) + "\n")
	}
	if pt.AdminUser != "" {
		s.WriteString(labelStyle.Render("Admin:   ") + valueStyle.Render(pt.AdminUser) + "\n")
	}

	return s.String()
}

// / docky-go logo
// / generated using: https://patorjk.com/software/taag/#p=display&f=Graffiti&t=docky-go%0A&x=none&v=4&h=4&w=80&we=false
var dockyLogo = []string{
	"    .___             __                                   ",
	"  __| _/____   ____ |  | _____.__.           ____   ____  ",
	" / __ |/  _ \\_/ ___\\|  |/ <   |  |  ______  / ___\\ /  _ \\ ",
	"/ /_/ (  <_> )  \\___|    < \\___  | /_____/ / /_/  >  <_> )",
	"\\____ |\\____/ \\___  >__|_ \\/ ____|         \\___  / \\____/ ",
	"     \\/           \\/     \\/\\/             /_____/          ",
}

func centeredLogo() string {
	var lines []string
	for _, l := range dockyLogo {
		lines = append(lines, strings.Repeat(" ", 0)+l)
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FDF500")).
		Bold(true).
		Render(strings.Join(lines, "\n"))
}
