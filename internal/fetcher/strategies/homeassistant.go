package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type HomeAssistantContainerInfo struct {
	model.BaseContainerInfo
	Version string
	Port    int
}

type HomeAssistantStrategy struct{}

func (s *HomeAssistantStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "homeassistant")
}

func (s *HomeAssistantStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &HomeAssistantContainerInfo{BaseContainerInfo: base}
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
		if ports, ok := inspect.NetworkSettings.Ports["8123/tcp"]; ok && len(ports) > 0 {
			if ports[0].HostPort != "" {
				fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
			}
		}
	}
	return info
}

func (h *HomeAssistantContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if h.Port > 0 {
		m["Port"] = fmt.Sprintf("%d", h.Port)
	}
	return m
}
