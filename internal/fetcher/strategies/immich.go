package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type ImmichContainerInfo struct {
	model.BaseContainerInfo
	Version   string
	DBHost    string
	DBPort    int
	RedisHost string
	RedisPort int
}

type ImmichStrategy struct{}

func (s *ImmichStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "immich")
}

func (s *ImmichStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &ImmichContainerInfo{BaseContainerInfo: base}
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
	}
	if err := dockerClient.ContainerInspect(ctx, id, &inspect); err == nil {
		envMap := model.ParseEnv(inspect.Config.Env)
		if v, ok := envMap["IMMICH_VERSION"]; ok {
			info.Version = v
		}
		if v, ok := envMap["DB_HOST"]; ok {
			info.DBHost = v
		}
		if v, ok := envMap["DB_PORT"]; ok {
			fmt.Sscanf(v, "%d", &info.DBPort)
		}
		if v, ok := envMap["REDIS_HOST"]; ok {
			info.RedisHost = v
		}
		if v, ok := envMap["REDIS_PORT"]; ok {
			fmt.Sscanf(v, "%d", &info.RedisPort)
		}
	}
	return info
}

func (i *ImmichContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if i.Version != "" {
		m["Version"] = i.Version
	}
	if i.DBHost != "" {
		m["DB Host"] = i.DBHost
	}
	if i.DBPort > 0 {
		m["DB Port"] = fmt.Sprintf("%d", i.DBPort)
	}
	if i.RedisHost != "" {
		m["Redis Host"] = i.RedisHost
	}
	if i.RedisPort > 0 {
		m["Redis Port"] = fmt.Sprintf("%d", i.RedisPort)
	}
	return m
}
