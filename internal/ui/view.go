package ui

import (
	"fmt"
	"sort"

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

	cols, _, _ := m.layoutSpec()
	width := m.termSize.Width
	if width <= 0 {
		width = 120
	}

	boxWidth := (width / cols) - 4
	columnContents := make([][]string, cols)
	columnHeights := make([]int, cols)

	for _, container := range m.items {
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
	for i := 0; i < cols; i++ {
		columnsRendered = append(columnsRendered, lipgloss.JoinVertical(lipgloss.Left, columnContents[i]...))
	}

	grid := lipgloss.JoinHorizontal(lipgloss.Top, columnsRendered...)

	footer := lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary)).Render("q-quit")
	return lipgloss.JoinVertical(lipgloss.Left, grid, footer)
}

func sortByHeight(boxes []string) []string {
	type pair struct {
		h int
		s string
	}
	var arr []pair
	for _, b := range boxes {
		arr = append(arr, pair{h: lipgloss.Height(b), s: b})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].h < arr[j].h })
	res := make([]string, len(arr))
	for i, p := range arr {
		res[i] = p.s
	}
	return res
}
