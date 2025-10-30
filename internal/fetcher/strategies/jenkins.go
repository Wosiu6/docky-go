package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type JenkinsContainerInfo struct {
	model.BaseContainerInfo
	Version   string
	AdminUser string
	Port      int
}

type JenkinsStrategy struct{}

func (s *JenkinsStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "jenkins")
}

func (s *JenkinsStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &JenkinsContainerInfo{BaseContainerInfo: base}
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
		if v, ok := envMap["JENKINS_VERSION"]; ok {
			info.Version = v
		}
		if v, ok := envMap["JENKINS_ADMIN_ID"]; ok {
			info.AdminUser = v
		}
		if ports, ok := inspect.NetworkSettings.Ports["8080/tcp"]; ok && len(ports) > 0 {
			if ports[0].HostPort != "" {
				fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
			}
		}
	}
	return info
}

func (j *JenkinsContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if j.Version != "" {
		m["Version"] = j.Version
	}
	if j.AdminUser != "" {
		m["Admin User"] = j.AdminUser
	}
	if j.Port > 0 {
		m["Port"] = fmt.Sprintf("%d", j.Port)
	}
	return m
}
