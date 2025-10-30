package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wosiu6/docky-go/internal/docker"
	"github.com/wosiu6/docky-go/internal/domain"
	"github.com/wosiu6/docky-go/internal/fetcher"
	"github.com/wosiu6/docky-go/internal/log"
	"github.com/wosiu6/docky-go/internal/orchestrator"
	"github.com/wosiu6/docky-go/internal/ui"
)

func main() {

	logger := log.New()

	dockerClient, err := docker.NewClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create docker client:", err)
		os.Exit(1)
	}

	pingCtx, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()
	if err := dockerClient.Ping(pingCtx); err != nil {
		if runtime.GOOS == "windows" {
			fmt.Fprintln(os.Stderr, "Cannot reach Docker. On Windows ensure Docker Desktop is running and named pipe \\ \\ . \\ pipe \\ docker_engine is available.")
		} else {
			fmt.Fprintln(os.Stderr, "Cannot reach Docker. Ensure the Docker daemon is running and /var/run/docker.sock is accessible.")
		}
		os.Exit(1)
	}

	dockerService := docker.NewService(dockerClient)
	f := fetcher.NewWithService(dockerService, dockerClient)
	adapter := fetcher.NewServiceAdapter(f)

	uiModel := ui.New(f)
	app := &uiAdapter{model: uiModel}

	orch := orchestrator.New(adapter, app, logger, 1*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := orch.Start(ctx); err != nil && err != context.Canceled {
		fmt.Fprintln(os.Stderr, "application error:", err)
		os.Exit(1)
	}
}

type uiAdapter struct {
	model   *ui.UiModel
	program *tea.Program
}

func (u *uiAdapter) SetData(containers []domain.Container) {
	legacy := make([]fetcher.ContainerInfo, 0, len(containers))
	for _, c := range containers {
		legacy = append(legacy, fetcher.ContainerInfo{
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
	u.model.SetItems(legacy)
	if u.program != nil {
		u.program.Send(ui.RefreshMsg{})
	}
}

func (u *uiAdapter) Run() error {
	u.program = tea.NewProgram(u.model)
	_, err := u.program.Run()
	return err
}

func (u *uiAdapter) NotifyRefresh() {
	if u.program != nil {
		u.program.Send(ui.RefreshMsg{})
	}
}
