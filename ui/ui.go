package ui

import (
	"context"
	"fmt"
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
	starting bool
}

func New(fetcher *fetcher.Fetcher) *UiModel { return &UiModel{fetcher: fetcher} }

func (m *UiModel) Init() tea.Cmd { return tickCmd() }

func tickCmd() tea.Cmd { return tea.Tick(time.Second*2, func(t time.Time) tea.Msg { return t }) }

func (m *UiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case time.Time:
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		items, err := m.fetcher.FetchAll(ctx)
		if err != nil {
			m.lastErr = err
		} else {
			m.items = items
			m.lastErr = nil
		}
		return m, tickCmd()
	}
	return m, nil
}

func (m *UiModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("Docker CLI Dashboard") + "\n\n"
	if m.lastErr != nil {
		header += lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("Error: "+m.lastErr.Error()) + "\n"
	}
	if m.starting {
		return header + "Loading...\n"
	} else if len(m.items) == 0 {
		return header + "No containers found.\n"
	}
	s := header + fmt.Sprintf("%-12s %-20s %-20s %-8s %-10s\n", "CONTAINER", "NAME", "IMAGE", "CPU%", "MEM")
	for _, it := range m.items {
		name := ""
		if len(it.Names) > 0 {
			name = strings.TrimPrefix(it.Names[0], "/")
		}
		s += fmt.Sprintf("%-12s %-20s %-20s %7.2f %-10d\n", it.ID[:12], name, it.Image, it.CPUPercent, it.Mem)
	}
	return s
}
