package fetcher

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wosiu6/docky-go/docker"
)

type BaseContainerInfo struct {
	ID         string
	Names      []string
	Image      string
	CPUPercent float64
	Mem        uint64
	Status     string
}

type PostgreSqlContainerInfo struct {
	BaseContainerInfo
	Port           int
	Database       string
	User           string
	SSLMode        string
	MaxConnections int
	PGData         string
}

type MinecraftContainerInfo struct {
	BaseContainerInfo
	Port          int
	Version       string
	ServerType    string
	Difficulty    string
	MaxPlayers    int
	OnlinePlayers int
}

type PortainerContainerInfo struct {
	BaseContainerInfo
	Port      int
	AdminUser string
	Edition   string
}

type ContainerType string

const (
	TypeGeneric    ContainerType = "generic"
	TypePostgreSQL ContainerType = "postgresql"
	TypeMinecraft  ContainerType = "minecraft"
	TypePortainer  ContainerType = "portainer"
)

type ContainerInfo struct {
	Type ContainerType
	BaseContainerInfo
	PostgreSql *PostgreSqlContainerInfo `json:",omitempty"`
	Minecraft  *MinecraftContainerInfo  `json:",omitempty"`
	Portainer  *PortainerContainerInfo  `json:",omitempty"`
}

type Fetcher struct {
	client *docker.DockerClient
	mu     sync.Mutex
	prev   map[string]StatsSnapshot
}

type StatsSnapshot struct {
	CPUTotal   uint64
	SystemCPU  uint64
	OnlineCPUs uint64
	Time       time.Time
}

func New(c *docker.DockerClient) *Fetcher {
	return &Fetcher{client: c, prev: make(map[string]StatsSnapshot)}
}

func (f *Fetcher) FetchAll(ctx context.Context) ([]ContainerInfo, error) {
	raw, err := f.client.ListContainers(ctx)
	if err != nil {
		return nil, err
	}

	type result struct {
		info ContainerInfo
		err  error
	}
	ch := make(chan result, len(raw))
	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup

	for _, r := range raw {
		id, _ := r["Id"].(string)
		namesIface, _ := r["Names"].([]interface{})
		image, _ := r["Image"].(string)
		state, _ := r["State"].(string)
		status, _ := r["Status"].(string)

		names := make([]string, 0, len(namesIface))
		for _, ni := range namesIface {
			if s, ok := ni.(string); ok {
				names = append(names, s)
			}
		}

		wg.Add(1)
		go func(id string, names []string, image string, state string, status string, rawContainer map[string]interface{}) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			var v struct {
				CPUStats struct {
					CPUUsage struct {
						TotalUsage uint64   `json:"total_usage"`
						Percpu     []uint64 `json:"percpu_usage"`
					} `json:"cpu_usage"`
					SystemCPUUsage uint64 `json:"system_cpu_usage"`
					OnlineCPUs     uint64 `json:"online_cpus"`
				} `json:"cpu_stats"`
				MemoryStats struct {
					Usage uint64 `json:"usage"`
				} `json:"memory_stats"`
			}

			if err := f.client.ContainerStats(ctx, id, &v); err != nil {
				ch <- result{
					info: ContainerInfo{
						Type: TypeGeneric,
						BaseContainerInfo: BaseContainerInfo{
							ID:     id,
							Names:  names,
							Image:  image,
							Status: state,
						},
					},
					err: nil,
				}
				return
			}

			snap := StatsSnapshot{
				CPUTotal:   v.CPUStats.CPUUsage.TotalUsage,
				SystemCPU:  v.CPUStats.SystemCPUUsage,
				OnlineCPUs: v.CPUStats.OnlineCPUs,
				Time:       time.Now(),
			}
			if snap.OnlineCPUs == 0 {
				snap.OnlineCPUs = uint64(len(v.CPUStats.CPUUsage.Percpu))
			}

			f.mu.Lock()
			prev, ok := f.prev[id]
			f.prev[id] = snap
			f.mu.Unlock()

			cpu := 0.0
			if ok {
				cpuDelta := float64(snap.CPUTotal - prev.CPUTotal)
				sysDelta := float64(snap.SystemCPU - prev.SystemCPU)
				if sysDelta > 0 && cpuDelta > 0 {
					cpu = (cpuDelta / sysDelta) * float64(snap.OnlineCPUs) * 100.0
				}
			}

			base := BaseContainerInfo{
				ID:         id,
				Names:      names,
				Image:      image,
				CPUPercent: cpu,
				Mem:        v.MemoryStats.Usage / 1024 / 1024,
				Status:     state,
			}

			containerType, specificInfo := f.detectContainerType(ctx, id, image, rawContainer, base)

			info := ContainerInfo{
				Type:              containerType,
				BaseContainerInfo: base,
			}

			switch containerType {
			case TypePostgreSQL:
				if pg, ok := specificInfo.(*PostgreSqlContainerInfo); ok {
					info.PostgreSql = pg
				}
			case TypeMinecraft:
				if mc, ok := specificInfo.(*MinecraftContainerInfo); ok {
					info.Minecraft = mc
				}
			case TypePortainer:
				if pt, ok := specificInfo.(*PortainerContainerInfo); ok {
					info.Portainer = pt
				}
			}

			ch <- result{info: info, err: nil}
		}(id, names, image, state, status, r)
	}

	wg.Wait()
	close(ch)

	var out []ContainerInfo
	for r := range ch {
		if r.err != nil {
			return nil, r.err
		}
		out = append(out, r.info)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].CPUPercent > out[j].CPUPercent })
	return out, nil
}

func (f *Fetcher) detectContainerType(ctx context.Context, id, image string, raw map[string]interface{}, base BaseContainerInfo) (ContainerType, interface{}) {
	imageLower := strings.ToLower(image)

	if strings.Contains(imageLower, "postgres") {
		return TypePostgreSQL, f.fetchPostgreSqlInfo(ctx, id, raw, base)
	}

	if strings.Contains(imageLower, "minecraft") || strings.Contains(imageLower, "itzg/minecraft") {
		return TypeMinecraft, f.fetchMinecraftInfo(ctx, id, raw, base)
	}

	if strings.Contains(imageLower, "portainer") {
		return TypePortainer, f.fetchPortainerInfo(ctx, id, raw, base)
	}

	return TypeGeneric, nil
}

func (f *Fetcher) fetchPostgreSqlInfo(ctx context.Context, id string, raw map[string]interface{}, base BaseContainerInfo) *PostgreSqlContainerInfo {
	info := &PostgreSqlContainerInfo{
		BaseContainerInfo: base,
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

	if err := f.client.ContainerInspect(ctx, id, &inspect); err == nil {
		for _, env := range inspect.Config.Env {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key, val := parts[0], parts[1]

			switch key {
			case "POSTGRES_DB":
				info.Database = val
			case "POSTGRES_USER":
				info.User = val
			case "POSTGRES_SSL_MODE":
				info.SSLMode = val
			case "PGDATA":
				info.PGData = val
			case "POSTGRES_MAX_CONNECTIONS":
				var maxConn int
				fmt.Sscanf(val, "%d", &maxConn)
				info.MaxConnections = maxConn
			}
		}

		if ports, ok := inspect.NetworkSettings.Ports["5432/tcp"]; ok && len(ports) > 0 {
			if ports[0].HostPort != "" {
				var port int
				fmt.Sscanf(ports[0].HostPort, "%d", &port)
				info.Port = port
			}
		}
	}

	return info
}

func (f *Fetcher) fetchMinecraftInfo(ctx context.Context, id string, raw map[string]interface{}, base BaseContainerInfo) *MinecraftContainerInfo {
	info := &MinecraftContainerInfo{
		BaseContainerInfo: base,
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

	if err := f.client.ContainerInspect(ctx, id, &inspect); err == nil {
		for _, env := range inspect.Config.Env {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key, val := parts[0], parts[1]

			switch key {
			case "VERSION":
				info.Version = val
			case "TYPE":
				info.ServerType = val
			case "DIFFICULTY":
				info.Difficulty = val
			case "MAX_PLAYERS":
				var maxPlayers int
				fmt.Sscanf(val, "%d", &maxPlayers)
				info.MaxPlayers = maxPlayers
			}
		}

		if ports, ok := inspect.NetworkSettings.Ports["25565/tcp"]; ok && len(ports) > 0 {
			if ports[0].HostPort != "" {
				var port int
				fmt.Sscanf(ports[0].HostPort, "%d", &port)
				info.Port = port
			}
		}
	}

	return info
}

func (f *Fetcher) fetchPortainerInfo(ctx context.Context, id string, raw map[string]interface{}, base BaseContainerInfo) *PortainerContainerInfo {
	info := &PortainerContainerInfo{
		BaseContainerInfo: base,
		Edition:           "Community",
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

	if err := f.client.ContainerInspect(ctx, id, &inspect); err == nil {
		for _, env := range inspect.Config.Env {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key, val := parts[0], parts[1]

			switch key {
			case "PORTAINER_ADMIN_USER":
				info.AdminUser = val
			}
		}

		if strings.Contains(strings.ToLower(base.Image), "portainer-ee") {
			info.Edition = "Business"
		}

		for portKey, ports := range inspect.NetworkSettings.Ports {
			if (strings.HasPrefix(portKey, "9000") || strings.HasPrefix(portKey, "9443")) && len(ports) > 0 {
				if ports[0].HostPort != "" {
					var port int
					fmt.Sscanf(ports[0].HostPort, "%d", &port)
					info.Port = port
					break
				}
			}
		}
	}

	return info
}
