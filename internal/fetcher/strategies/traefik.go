package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type TraefikContainerInfo struct {
	model.BaseContainerInfo
	Version     string
	Entrypoints string
	Dashboard   bool
	Ports       []int
}

func (t *TraefikContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if t.Version != "" {
		m["Version"] = t.Version
	}
	if t.Entrypoints != "" {
		m["Entrypoints"] = t.Entrypoints
	}
	if t.Dashboard {
		m["Dashboard"] = "Enabled"
	}
	if len(t.Ports) > 0 {
		var ps []string
		for _, p := range t.Ports {
			ps = append(ps, fmt.Sprintf("%d", p))
		}
		m["Ports"] = strings.Join(ps, ", ")
	}
	return m
}

type TraefikStrategy struct{}

func (s *TraefikStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "traefik")
}

func (s *TraefikStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &TraefikContainerInfo{BaseContainerInfo: base}
	dockerClient, ok := client.(interface {
		ContainerInspect(context.Context, string, interface{}) error
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
		if v, ok := envMap["TRAEFIK_VERSION"]; ok {
			info.Version = v
		}
		if v, ok := envMap["TRAEFIK_ENTRYPOINTS"]; ok {
			info.Entrypoints = v
		}
		if v, ok := envMap["TRAEFIK_DASHBOARD"]; ok {
			info.Dashboard = v == "true"
		}
		for _, ports := range inspect.NetworkSettings.Ports {
			if len(ports) > 0 && ports[0].HostPort != "" {
				var port int
				fmt.Sscanf(ports[0].HostPort, "%d", &port)
				info.Ports = append(info.Ports, port)
			}
		}
	}
	return info
}
