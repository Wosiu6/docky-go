package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type RabbitMQContainerInfo struct {
	model.BaseContainerInfo
	Version string
	User    string
	Port    int
}

type RabbitMQStrategy struct{}

func (s *RabbitMQStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "rabbitmq")
}

func (s *RabbitMQStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &RabbitMQContainerInfo{BaseContainerInfo: base}
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
		if v, ok := envMap["RABBITMQ_DEFAULT_USER"]; ok {
			info.User = v
		}
		if v, ok := envMap["RABBITMQ_VERSION"]; ok {
			info.Version = v
		}
		if ports, ok := inspect.NetworkSettings.Ports["5672/tcp"]; ok && len(ports) > 0 {
			if ports[0].HostPort != "" {
				fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
			}
		}
	}
	return info
}

func (r *RabbitMQContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if r.Version != "" {
		m["Version"] = r.Version
	}
	if r.User != "" {
		m["User"] = r.User
	}
	if r.Port > 0 {
		m["Port"] = fmt.Sprintf("%d", r.Port)
	}
	return m
}
