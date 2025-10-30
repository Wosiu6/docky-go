package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type MySQLContainerInfo struct {
	model.BaseContainerInfo
	Version  string
	User     string
	Database string
	Port     int
}

type MySQLStrategy struct{}

func (s *MySQLStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "mysql")
}

func (s *MySQLStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &MySQLContainerInfo{BaseContainerInfo: base}
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
		if v, ok := envMap["MYSQL_VERSION"]; ok {
			info.Version = v
		}
		if v, ok := envMap["MYSQL_USER"]; ok {
			info.User = v
		}
		if v, ok := envMap["MYSQL_DATABASE"]; ok {
			info.Database = v
		}
		if ports, ok := inspect.NetworkSettings.Ports["3306/tcp"]; ok && len(ports) > 0 {
			if ports[0].HostPort != "" {
				fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
			}
		}
	}
	return info
}

func (m *MySQLContainerInfo) DetailFields() map[string]string {
	mOut := map[string]string{}
	if m.Version != "" {
		mOut["Version"] = m.Version
	}
	if m.User != "" {
		mOut["User"] = m.User
	}
	if m.Database != "" {
		mOut["Database"] = m.Database
	}
	if m.Port > 0 {
		mOut["Port"] = fmt.Sprintf("%d", m.Port)
	}
	return mOut
}
