package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type MongoDBContainerInfo struct {
	model.BaseContainerInfo
	Version  string
	Port     int
	Database string
}

type MongoDBStrategy struct{}

func (s *MongoDBStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "mongo")
}

func (s *MongoDBStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &MongoDBContainerInfo{BaseContainerInfo: base}
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
		if v, ok := envMap["MONGO_INITDB_DATABASE"]; ok {
			info.Database = v
		}
		if ports, ok := inspect.NetworkSettings.Ports["27017/tcp"]; ok && len(ports) > 0 {
			if ports[0].HostPort != "" {
				fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
			}
		}
	}
	return info
}

func (m *MongoDBContainerInfo) DetailFields() map[string]string {
	mOut := map[string]string{}
	if m.Database != "" {
		mOut["Database"] = m.Database
	}
	if m.Port > 0 {
		mOut["Port"] = fmt.Sprintf("%d", m.Port)
	}
	return mOut
}
