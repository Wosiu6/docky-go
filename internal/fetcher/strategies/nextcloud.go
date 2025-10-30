package strategies

import (
	"context"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type NextcloudContainerInfo struct {
	model.BaseContainerInfo
	Version   string
	AdminUser string
	DBHost    string
	DBName    string
}

func (n *NextcloudContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if n.Version != "" {
		m["Version"] = n.Version
	}
	if n.AdminUser != "" {
		m["Admin User"] = n.AdminUser
	}
	if n.DBHost != "" {
		m["DB Host"] = n.DBHost
	}
	if n.DBName != "" {
		m["DB Name"] = n.DBName
	}
	return m
}

type NextcloudStrategy struct{}

func (s *NextcloudStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "nextcloud")
}

func (s *NextcloudStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &NextcloudContainerInfo{BaseContainerInfo: base}
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
	}
	if err := dockerClient.ContainerInspect(ctx, id, &inspect); err == nil {
		envMap := model.ParseEnv(inspect.Config.Env)
		if v, ok := envMap["NEXTCLOUD_VERSION"]; ok {
			info.Version = v
		}
		if v, ok := envMap["NEXTCLOUD_ADMIN_USER"]; ok {
			info.AdminUser = v
		}
		if v, ok := envMap["MYSQL_HOST"]; ok {
			info.DBHost = v
		}
		if v, ok := envMap["MYSQL_DATABASE"]; ok {
			info.DBName = v
		}
	}
	return info
}
