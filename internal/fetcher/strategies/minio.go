package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type MinioContainerInfo struct {
	model.BaseContainerInfo
	Version     string
	AccessKey   string
	SecretKey   string
	ConsolePort int
}

type MinioStrategy struct{}

func (s *MinioStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "minio")
}

func (s *MinioStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &MinioContainerInfo{BaseContainerInfo: base}
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
		if v, ok := envMap["MINIO_ROOT_USER"]; ok {
			info.AccessKey = v
		}
		if v, ok := envMap["MINIO_ROOT_PASSWORD"]; ok {
			info.SecretKey = v
		}
		if ports, ok := inspect.NetworkSettings.Ports["9001/tcp"]; ok && len(ports) > 0 {
			if ports[0].HostPort != "" {
				fmt.Sscanf(ports[0].HostPort, "%d", &info.ConsolePort)
			}
		}
	}
	return info
}

func (m *MinioContainerInfo) DetailFields() map[string]string {
	mOut := map[string]string{}
	if m.AccessKey != "" {
		mOut["Access Key"] = m.AccessKey
	}
	if m.SecretKey != "" {
		mOut["Secret Key"] = m.SecretKey
	}
	if m.ConsolePort > 0 {
		mOut["Console Port"] = fmt.Sprintf("%d", m.ConsolePort)
	}
	return mOut
}
