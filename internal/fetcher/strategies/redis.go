package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type RedisContainerInfo struct {
	model.BaseContainerInfo
	Port     int
	Password string
}

func (r *RedisContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if r.Port > 0 {
		m["Port"] = fmt.Sprintf("%d", r.Port)
	}
	if r.Password != "" {
		m["Password"] = r.Password
	}
	return m
}

type RedisStrategy struct{}

func (s *RedisStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "redis")
}

func (s *RedisStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &RedisContainerInfo{BaseContainerInfo: base}
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
		if v, ok := envMap["REDIS_PASSWORD"]; ok {
			info.Password = v
		}
		if ports, ok := inspect.NetworkSettings.Ports["6379/tcp"]; ok && len(ports) > 0 {
			fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
		}
	}
	return info
}
