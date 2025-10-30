package strategies

import (
	"context"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type OwnCloudContainerInfo struct {
	model.BaseContainerInfo
	Version   string
	AdminUser string
	AdminPass string
	DBHost    string
	DBName    string
}

type OwnCloudStrategy struct{}

func (s *OwnCloudStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "owncloud")
}

func (s *OwnCloudStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &OwnCloudContainerInfo{BaseContainerInfo: base}
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
		if v, ok := envMap["OWNCLOUD_VERSION"]; ok {
			info.Version = v
		}
		if v, ok := envMap["OWNCLOUD_ADMIN_USERNAME"]; ok {
			info.AdminUser = v
		}
		if v, ok := envMap["OWNCLOUD_ADMIN_PASSWORD"]; ok {
			info.AdminPass = v
		}
		if v, ok := envMap["OWNCLOUD_DB_HOST"]; ok {
			info.DBHost = v
		}
		if v, ok := envMap["OWNCLOUD_DB_NAME"]; ok {
			info.DBName = v
		}
	}
	return info
}

func (o *OwnCloudContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if o.Version != "" {
		m["Version"] = o.Version
	}
	if o.DBHost != "" {
		m["DB Host"] = o.DBHost
	}
	if o.DBName != "" {
		m["DB Name"] = o.DBName
	}
	return m
}
