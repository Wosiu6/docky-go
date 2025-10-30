package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type WordPressContainerInfo struct {
	model.BaseContainerInfo
	Version string
	DBHost  string
	DBName  string
	Port    int
}

func (w *WordPressContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if w.Version != "" {
		m["Version"] = w.Version
	}
	if w.DBHost != "" {
		m["DB Host"] = w.DBHost
	}
	if w.DBName != "" {
		m["DB Name"] = w.DBName
	}
	if w.Port > 0 {
		m["Port"] = fmt.Sprintf("%d", w.Port)
	}
	return m
}

type WordPressStrategy struct{}

func (s *WordPressStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "wordpress")
}

func (s *WordPressStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &WordPressContainerInfo{BaseContainerInfo: base}
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
		if v, ok := envMap["WORDPRESS_VERSION"]; ok {
			info.Version = v
		}
		if v, ok := envMap["WORDPRESS_DB_HOST"]; ok {
			info.DBHost = v
		}
		if v, ok := envMap["WORDPRESS_DB_NAME"]; ok {
			info.DBName = v
		}
		if ports, ok := inspect.NetworkSettings.Ports["80/tcp"]; ok && len(ports) > 0 {
			fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
		}
	}
	return info
}
