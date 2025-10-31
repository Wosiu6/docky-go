package ui

import (
	"fmt"
	"strings"

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

	cols, _, perPage := m.layoutSpec()
	if perPage < cols {
		perPage = cols
	}
	width := m.termSize.Width
	if width <= 0 {
		width = 120
	}

	boxWidth := (width / cols) - 4
	columnContents := make([][]string, cols)
	columnHeights := make([]int, cols)

	start := m.page * perPage
	if start >= len(m.items) {
		start = 0
		m.page = 0
	}
	end := start + perPage
	if end > len(m.items) {
		end = len(m.items)
	}
	visible := m.items[start:end]

	for _, container := range visible {
		box := m.renderContainer(container, boxWidth, 0)
		minIdx := 0
		minHeight := columnHeights[0]
		for i := 1; i < cols; i++ {
			if columnHeights[i] < minHeight {
				minIdx = i
				minHeight = columnHeights[i]
			}
		}
		columnContents[minIdx] = append(columnContents[minIdx], box)
		columnHeights[minIdx] += lipgloss.Height(box) + 1
	}

	var columnsRendered []string
	for i := range cols {
		columnsRendered = append(columnsRendered, lipgloss.JoinVertical(lipgloss.Left, columnContents[i]...))
	}

	grid := lipgloss.JoinHorizontal(lipgloss.Top, columnsRendered...)

	footer := m.renderFooter()
	return lipgloss.JoinVertical(lipgloss.Left, grid, footer)
}

func (m *UiModel) renderFooter() string {
	const sep = " │ "

	quit := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorDanger)).
		Render("q quit")

	if m.totalPages() <= 1 {
		return lipgloss.NewStyle().
			Width(m.termSize.Width).
			Render(quit)
	}

	navStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextDim))

	leftNav := navStyle.Render("h/← prev")
	rightNav := navStyle.Render("l/→ next")

	pageInfo := fmt.Sprintf(" %d / %d ", m.page+1, m.totalPages())
	page := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Bold(true).
		Render(pageInfo)

	navBlock := lipgloss.JoinHorizontal(lipgloss.Center, leftNav, sep, page, sep, rightNav)

	available := m.termSize.Width -
		lipgloss.Width(navBlock) -
		lipgloss.Width(quit) -
		lipgloss.Width(sep)

	middle := strings.Repeat(" ", max(0, available))

	footer := lipgloss.JoinHorizontal(lipgloss.Top,
		navBlock,
		middle,
		sep,
		quit,
	)

	return lipgloss.NewStyle().
		Width(m.termSize.Width).
		Render(footer)
}
