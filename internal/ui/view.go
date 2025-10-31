package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m *UiModel) View() string {
	if m.lastErr != nil {
		return errorStyle.Render(fmt.Sprintf("\u274c Error: %v", m.lastErr))
	}
	if m.loading {
		return logo() + "\n\nLoading containers..."
	}
	if len(m.items) == 0 {
		renderString := logo() + "\n\n\U0001F433 No containers found. Waiting for containers..."
		return emptyStyle.Render(renderString)
	}
	cols, rows, maxItems := m.layoutSpec()
	width := m.termSize.Width
	height := m.termSize.Height - 2
	if width <= 0 {
		width = 120
	}
	if height <= 0 {
		height = 30
	}
	boxWidth := (width / cols) - 4
	boxHeight := (height / rows) - 6
	start := m.page * maxItems
	end := start + maxItems
	if start >= len(m.items) {
		start = 0
		m.page = 0
		end = maxItems
	}
	if end > len(m.items) {
		end = len(m.items)
	}
	visible := m.items[start:end]
	var grid []string
	var currentRow []string
	for i, container := range visible {
		box := m.renderContainer(container, boxWidth, boxHeight)
		currentRow = append(currentRow, box)
		if len(currentRow) == cols || i == len(visible)-1 {
			grid = append(grid, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
			currentRow = []string{}
		}
	}
	pages := m.totalPages()
	footer := ""
	if pages > 1 {
		footer = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary)).Render(fmt.Sprintf("Page %d/%d (←h/l→/q-quit)", m.page+1, pages))
	}
	content := lipgloss.JoinVertical(lipgloss.Left, append(grid, footer)...)
	return content
}
