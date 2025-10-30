package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type PortainerContainerInfo struct {
	model.BaseContainerInfo
	Port      int
	AdminUser string
	Edition   string
}

func (pt *PortainerContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if pt.Port > 0 {
		m["Port"] = fmt.Sprintf("%d", pt.Port)
	}
	if pt.Edition != "" {
		m["Edition"] = pt.Edition
	}
	if pt.AdminUser != "" {
		m["Admin"] = pt.AdminUser
	}
	return m
}

type PortainerStrategy struct{}

func (s *PortainerStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "portainer")
}

func (s *PortainerStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &PortainerContainerInfo{BaseContainerInfo: base, Edition: "Community"}
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
		if v, ok := envMap["PORTAINER_ADMIN_USER"]; ok {
			info.AdminUser = v
		}
		if strings.Contains(strings.ToLower(base.Image), "portainer-ee") {
			info.Edition = "Business"
		}
		for portKey, ports := range inspect.NetworkSettings.Ports {
			if (strings.HasPrefix(portKey, "9000") || strings.HasPrefix(portKey, "9443")) && len(ports) > 0 {
				fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
				break
			}
		}
	}
	return info
}
