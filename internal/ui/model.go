package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wosiu6/docky-go/internal/fetcher"
)

type FetcherInterface interface {
	FetchAll(ctx context.Context) ([]fetcher.ContainerInfo, error)
}

type UiModel struct {
	fetcher  FetcherInterface
	lastErr  error
	items    []fetcher.ContainerInfo
	loading  bool
	termSize tea.WindowSizeMsg
	page     int
}

type RefreshMsg struct{}

func New(fetcher FetcherInterface) *UiModel { return &UiModel{fetcher: fetcher, loading: true} }
func (m *UiModel) SetItems(items []fetcher.ContainerInfo) {
	m.items = items
	m.loading = false
}

func (m *UiModel) Init() tea.Cmd { return tea.ClearScreen }

func (m *UiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l":
			m.nextPage()
			return m, nil
		case "left", "h":
			m.prevPage()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.termSize = msg
		return m, nil
	case RefreshMsg:
		return m, nil
	}
	return m, nil
}

func (m *UiModel) nextPage() {
	maxPages := m.totalPages()
	if m.page < maxPages-1 {
		m.page++
	}
}

func (m *UiModel) prevPage() {
	if m.page > 0 {
		m.page--
	}
}

func (m *UiModel) totalPages() int {
	if len(m.items) == 0 {
		return 1
	}
	_, _, perPage := m.layoutSpec()
	if perPage <= 0 {
		return 1
	}
	pages := (len(m.items) + perPage - 1) / perPage
	if pages < 1 {
		pages = 1
	}
	return pages
}

func (m *UiModel) layoutSpec() (int, int, int) {
	width := m.termSize.Width
	height := m.termSize.Height - 2
	if width <= 0 {
		width = 120
	}
	if height <= 0 {
		height = 30
	}
	minBoxH := 10
	minBoxW := 30
	cols := 1
	if width >= minBoxW*2+8 {
		cols = 2
	}
	if width >= minBoxW*3+12 {
		cols = 3
	}
	rows := max(height/(minBoxH+2), 1)
	if rows > 3 {
		rows = 3
	}
	maxItems := cols * rows
	return cols, rows, maxItems
}
