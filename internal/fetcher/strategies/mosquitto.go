package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type MosquittoContainerInfo struct {
	model.BaseContainerInfo
	Version string
	Port    int
}

type MosquittoStrategy struct{}

func (s *MosquittoStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "mosquitto")
}

func (s *MosquittoStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &MosquittoContainerInfo{BaseContainerInfo: base}
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
		if ports, ok := inspect.NetworkSettings.Ports["1883/tcp"]; ok && len(ports) > 0 {
			if ports[0].HostPort != "" {
				fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
			}
		}
	}
	return info
}

func (m *MosquittoContainerInfo) DetailFields() map[string]string {
	mOut := map[string]string{}
	if m.Port > 0 {
		mOut["Port"] = fmt.Sprintf("%d", m.Port)
	}
	return mOut
}
