package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wosiu6/docky-go/docker"
	"github.com/wosiu6/docky-go/fetcher"
	"github.com/wosiu6/docky-go/ui"
)

func main() {
	dockerClient, err := docker.NewClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create docker client:", err)
		os.Exit(1)
	}

	// quick check: ping
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", dockerClient.GetUrl()+"/_ping", nil)
	dockerClientHttp := dockerClient.GetHttpClient()
	if resp, err := dockerClientHttp.Do(req); err != nil || resp.StatusCode >= 400 {
		if runtime.GOOS == "windows" {
			fmt.Fprintln(os.Stderr, "Cannot reach Docker. On Windows ensure Docker Desktop is running and named pipe \\\\.\\pipe\\docker_engine is available.")
		} else {
			fmt.Fprintln(os.Stderr, "Cannot reach Docker. Ensure the Docker daemon is running and /var/run/docker.sock is accessible.")
		}
		if resp != nil && resp.Body != nil {
			io.Copy(io.Discard, resp.Body)
		}
		os.Exit(1)
	}

	fetcher := fetcher.New(dockerClient)
	uiModel := ui.New(fetcher)
	program := tea.NewProgram(uiModel)
	if _, err := program.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "TUI error:", err)
		os.Exit(1)
	}
}
