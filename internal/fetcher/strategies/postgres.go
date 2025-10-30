package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/wosiu6/docky-go/internal/model"
)

type PostgreSqlContainerInfo struct {
	model.BaseContainerInfo
	Port           int
	Database       string
	User           string
	SSLMode        string
	MaxConnections int
	PGData         string
}

func (pg *PostgreSqlContainerInfo) DetailFields() map[string]string {
	m := map[string]string{}
	if pg.Port > 0 {
		m["Port"] = fmt.Sprintf("%d", pg.Port)
	}
	if pg.Database != "" {
		m["Database"] = pg.Database
	}
	if pg.User != "" {
		m["User"] = pg.User
	}
	if pg.SSLMode != "" {
		m["SSL Mode"] = pg.SSLMode
	}
	if pg.MaxConnections > 0 {
		m["Max Conn"] = fmt.Sprintf("%d", pg.MaxConnections)
	}
	if pg.PGData != "" {
		m["Volume"] = pg.PGData
	}
	return m
}

type PostgreSqlStrategy struct{}

func (s *PostgreSqlStrategy) Match(image string) bool {
	return strings.Contains(strings.ToLower(image), "postgres")
}

func (s *PostgreSqlStrategy) Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{} {
	info := &PostgreSqlContainerInfo{BaseContainerInfo: base}
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
		if v, ok := envMap["POSTGRES_DB"]; ok {
			info.Database = v
		}
		if v, ok := envMap["POSTGRES_USER"]; ok {
			info.User = v
		}
		if v, ok := envMap["POSTGRES_SSL_MODE"]; ok {
			info.SSLMode = v
		}
		if v, ok := envMap["PGDATA"]; ok {
			info.PGData = v
		}
		if v, ok := envMap["POSTGRES_MAX_CONNECTIONS"]; ok {
			fmt.Sscanf(v, "%d", &info.MaxConnections)
		}
		if ports, ok := inspect.NetworkSettings.Ports["5432/tcp"]; ok && len(ports) > 0 {
			fmt.Sscanf(ports[0].HostPort, "%d", &info.Port)
		}
	}
	return info
}
