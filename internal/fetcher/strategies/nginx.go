package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type NginxContainerInfo struct {
	model.BaseContainerInfo
	Ports []int
}

func (n *NginxContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if len(n.Ports) > 0 {
		var ps []string
		for _, p := range n.Ports {
			ps = append(ps, fmt.Sprintf("%d", p))
		}
		m["Ports"] = strings.Join(ps, ", ")
	}
	return m
}

type NginxStrategy struct{}

func (s *NginxStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "nginx")
}

func (s *NginxStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &NginxContainerInfo{BaseContainerInfo: base}
	dockerClient, ok := client.(interface {
		ContainerInspect(context.Context, string, interface{}) error
	})
	if !ok {
		return info
	}
	var inspect struct {
		NetworkSettings struct {
			Ports map[string][]struct {
				HostPort string `json:"HostPort"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}
	if err := dockerClient.ContainerInspect(ctx, id, &inspect); err == nil {
		for _, ports := range inspect.NetworkSettings.Ports {
			if len(ports) > 0 && ports[0].HostPort != "" {
				var port int
				fmt.Sscanf(ports[0].HostPort, "%d", &port)
				info.Ports = append(info.Ports, port)
			}
		}
	}
	return info
}
