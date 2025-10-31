package main

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/wosiu6/docky-go/internal/docker"
	"github.com/wosiu6/docky-go/internal/fetcher"
	"github.com/wosiu6/docky-go/internal/log"
	"github.com/wosiu6/docky-go/internal/orchestrator"
	"github.com/wosiu6/docky-go/internal/ui"
)

func main() {

	logger := log.New()

	dockerClient, err := docker.NewClient()
	if err != nil {
		logger.Error("failed to create docker client", "error", err)
		os.Exit(1)
	}

	pingCtx, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()
	if err := dockerClient.Ping(pingCtx); err != nil {
		if runtime.GOOS == "windows" {
			logger.Error("Cannot reach Docker. On Windows ensure Docker Desktop is running and named pipe \\ \\ . \\ pipe \\ docker_engine is available.", "error", err)
		} else {
			logger.Error("Cannot reach Docker. Ensure the Docker daemon is running and /var/run/docker.sock is accessible. Ensure you have permission to access the Docker socket/are a part of the docker group.", "error", err)
		}
		os.Exit(1)
	}

	dockerService := docker.NewService(dockerClient)
	containerFetcher := fetcher.NewWithService(dockerService, dockerClient)
	serviceAdapter := fetcher.NewServiceAdapter(containerFetcher)

	uiModel := ui.New(containerFetcher)
	uiAdapter := ui.NewAdapter(uiModel)

	orchestrator := orchestrator.New(serviceAdapter, uiAdapter, logger, time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := orchestrator.Start(ctx); err != nil && err != context.Canceled {
		logger.Error("application error", "error", err)
		os.Exit(1)
	}
}
