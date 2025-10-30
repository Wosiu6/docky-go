package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wosiu6/docky-go/internal/domain"
	"github.com/wosiu6/docky-go/internal/fetcher"
)

type UiAdapter struct {
	model   *UiModel
	program *tea.Program
}

func NewAdapter(model *UiModel) *UiAdapter {
	return &UiAdapter{model: model}
}

func (u *UiAdapter) SetData(containers []domain.Container) {
	containerInfoList := make([]fetcher.ContainerInfo, 0, len(containers))
	for _, c := range containers {
		containerInfoList = append(containerInfoList, fetcher.ContainerInfo{
			Type: c.Type,
			BaseContainerInfo: fetcher.BaseContainerInfo{
				ID:         c.ID,
				Names:      c.Names,
				Image:      c.Image,
				CPUPercent: c.CPUPercent,
				Mem:        c.MemoryMB,
				Status:     c.Status,
			},
			Specific: c.Details,
		})
	}
	u.model.SetItems(containerInfoList)
	if u.program != nil {
		u.program.Send(RefreshMsg{})
	}
}

func (u *UiAdapter) Run() error {
	u.program = tea.NewProgram(u.model)
	_, err := u.program.Run()
	return err
}

func (u *UiAdapter) NotifyRefresh() {
	if u.program != nil {
		u.program.Send(RefreshMsg{})
	}
}
