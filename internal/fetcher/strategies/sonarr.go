package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type SonarrContainerInfo struct {
	model.BaseContainerInfo
	Version string
	Port    int
}

type SonarrStrategy struct{}

func (s *SonarrStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "sonarr")
}

func (s *SonarrStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &SonarrContainerInfo{BaseContainerInfo: base}
	dockerClient, ok := client.(interface {
		ContainerInspect(ctx context.Context, id string, v interface{}) error
	})
	if !ok {
		return info
	}

	var inspect struct {
		Config struct {
			Env []string `json:"Env"`
		} `json:"Config"`
		NetworkSettings struct {
			Ports map[string][]struct {
				HostPort string `json:"HostPort"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}
	if err := dockerClient.ContainerInspect(ctx, id, &inspect); err == nil {
		if ports, ok := inspect.NetworkSettings.Ports["8989/tcp"]; ok && len(ports) > 0 {
			if ports[0].HostPort != "" {
				fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
			}
		}
	}
	return info
}

func (s *SonarrContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if s.Port > 0 {
		m["Port"] = fmt.Sprintf("%d", s.Port)
	}
	return m
}
