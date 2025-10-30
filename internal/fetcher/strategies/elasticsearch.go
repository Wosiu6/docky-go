package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type ElasticsearchContainerInfo struct {
	model.BaseContainerInfo
	Version string
	Port    int
}

type ElasticsearchStrategy struct{}

func (s *ElasticsearchStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "elasticsearch")
}

func (s *ElasticsearchStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &ElasticsearchContainerInfo{BaseContainerInfo: base}
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
		envMap := model.ParseEnv(inspect.Config.Env)
		if v, ok := envMap["ELASTIC_VERSION"]; ok {
			info.Version = v
		}
		if ports, ok := inspect.NetworkSettings.Ports["9200/tcp"]; ok && len(ports) > 0 {
			if ports[0].HostPort != "" {
				fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
			}
		}
	}
	return info
}

func (e *ElasticsearchContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if e.Version != "" {
		m["Version"] = e.Version
	}
	if e.Port > 0 {
		m["Port"] = fmt.Sprintf("%d", e.Port)
	}
	return m
}
