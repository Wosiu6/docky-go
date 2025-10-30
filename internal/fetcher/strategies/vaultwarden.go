package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type VaultwardenContainerInfo struct {
	model.BaseContainerInfo
	Version    string
	AdminToken string
	Port       int
}

func (v *VaultwardenContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if v.Version != "" {
		m["Version"] = v.Version
	}
	if v.AdminToken != "" {
		m["Admin Token"] = v.AdminToken
	}
	if v.Port > 0 {
		m["Port"] = fmt.Sprintf("%d", v.Port)
	}
	return m
}

type VaultwardenStrategy struct{}

func (s *VaultwardenStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "vaultwarden")
}

func (s *VaultwardenStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &VaultwardenContainerInfo{BaseContainerInfo: base}
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
		if v, ok := envMap["ADMIN_TOKEN"]; ok {
			info.AdminToken = v
		}
		if v, ok := envMap["VAULTWARDEN_VERSION"]; ok {
			info.Version = v
		}
		if ports, ok := inspect.NetworkSettings.Ports["80/tcp"]; ok && len(ports) > 0 {
			fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
		}
	}
	return info
}
