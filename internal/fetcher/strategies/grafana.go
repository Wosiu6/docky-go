package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type GrafanaContainerInfo struct {
	model.BaseContainerInfo
	Version   string
	AdminUser string
	Port      int
}

type GrafanaStrategy struct{}

func (s *GrafanaStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "grafana")
}

func (s *GrafanaStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &GrafanaContainerInfo{BaseContainerInfo: base}
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
		if v, ok := envMap["GF_SECURITY_ADMIN_USER"]; ok {
			info.AdminUser = v
		}
		if v, ok := envMap["GF_VERSION"]; ok {
			info.Version = v
		}
		if ports, ok := inspect.NetworkSettings.Ports["3000/tcp"]; ok && len(ports) > 0 {
			if ports[0].HostPort != "" {
				fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
			}
		}
	}
	return info
}

func (g *GrafanaContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if g.Version != "" {
		m["Version"] = g.Version
	}
	if g.AdminUser != "" {
		m["Admin User"] = g.AdminUser
	}
	if g.Port > 0 {
		m["Port"] = fmt.Sprintf("%d", g.Port)
	}
	return m
}
