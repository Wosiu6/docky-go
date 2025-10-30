package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type MinecraftContainerInfo struct {
	model.BaseContainerInfo
	Port          int
	Version       string
	ServerType    string
	Difficulty    string
	MaxPlayers    int
	OnlinePlayers int
}

func (mc *MinecraftContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if mc.Port > 0 {
		m["Port"] = fmt.Sprintf("%d", mc.Port)
	}
	if mc.Version != "" {
		m["Version"] = mc.Version
	}
	if mc.ServerType != "" {
		m["Type"] = mc.ServerType
	}
	if mc.Difficulty != "" {
		m["Difficulty"] = mc.Difficulty
	}
	if mc.MaxPlayers > 0 {
		m["Players"] = fmt.Sprintf("%d/%d", mc.OnlinePlayers, mc.MaxPlayers)
	}
	return m
}

type MinecraftStrategy struct{}

func (s *MinecraftStrategy) Match(image string) bool {
	img := strings.ToLower(image)
	return strings.Contains(img, "minecraft") || strings.Contains(img, "itzg/minecraft")
}

func (s *MinecraftStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &MinecraftContainerInfo{BaseContainerInfo: base}
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
		if v, ok := envMap["VERSION"]; ok {
			info.Version = v
		}
		if v, ok := envMap["TYPE"]; ok {
			info.ServerType = v
		}
		if v, ok := envMap["DIFFICULTY"]; ok {
			info.Difficulty = v
		}
		if v, ok := envMap["MAX_PLAYERS"]; ok {
			fmt.Sscanf(v, "%d", &info.MaxPlayers)
		}
		if ports, ok := inspect.NetworkSettings.Ports["25565/tcp"]; ok && len(ports) > 0 {
			fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
		}
	}
	return info
}
