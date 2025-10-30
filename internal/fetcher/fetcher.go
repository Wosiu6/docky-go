package fetcher

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wosiu6/docky-go/internal/docker"
	"github.com/wosiu6/docky-go/internal/domain"
	"github.com/wosiu6/docky-go/internal/fetcher/strategies"
	"github.com/wosiu6/docky-go/internal/model"
)

type BaseContainerInfo = model.BaseContainerInfo

type DetailProvider interface{ DetailFields() map[string]string }

type ContainerInfo struct {
	Type domain.ContainerType
	BaseContainerInfo
	Specific DetailProvider
}

type Fetcher struct {
	client  docker.DockerClient
	service docker.Service
	mu      sync.Mutex
	prev    map[string]StatsSnapshot
	entries []strategies.StrategyEntry
	cfg     FetcherConfig
}

type FetcherConfig struct {
	Concurrency int
	SortByCPU   bool
}

func defaultConfig() FetcherConfig { return FetcherConfig{Concurrency: 8, SortByCPU: false} }

type StatsSnapshot struct {
	CPUTotal   uint64
	SystemCPU  uint64
	OnlineCPUs uint64
	Time       time.Time
}

func New(c docker.DockerClient) *Fetcher { return NewWithConfig(c, defaultConfig()) }

func NewWithConfig(c docker.DockerClient, cfg FetcherConfig) *Fetcher {
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 4
	}
	return &Fetcher{client: c, prev: make(map[string]StatsSnapshot), entries: strategies.Registry(), cfg: cfg}
}

func NewWithService(s docker.Service, raw docker.DockerClient) *Fetcher {
	f := New(raw)
	f.service = s
	return f
}

func (f *Fetcher) FetchAll(ctx context.Context) ([]ContainerInfo, error) {
	var raw []map[string]interface{}
	var err error
	if f.service != nil {
		raw, err = f.service.Containers(ctx)
	} else {
		raw, err = f.client.ListContainers(ctx)
	}
	if err != nil {
		return nil, err
	}
	type result struct {
		info ContainerInfo
		err  error
	}
	ch := make(chan result, len(raw))
	sem := make(chan struct{}, f.cfg.Concurrency)
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
				ch <- result{info: ContainerInfo{Type: domain.ContainerTypeGeneric, BaseContainerInfo: BaseContainerInfo{ID: id, Names: names, Image: image, Status: state}}, err: nil}
				return
			}
			snap := StatsSnapshot{CPUTotal: v.CPUStats.CPUUsage.TotalUsage, SystemCPU: v.CPUStats.SystemCPUUsage, OnlineCPUs: v.CPUStats.OnlineCPUs, Time: time.Now()}
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
			base := model.BaseContainerInfo{ID: id, Names: names, Image: image, CPUPercent: cpu, Mem: v.MemoryStats.Usage / 1024 / 1024, Status: state}
			var matchedType domain.ContainerType = domain.ContainerTypeGeneric
			var specific DetailProvider
			for _, entry := range f.entries {
				if entry.Strategy.Match(image) {
					matchedType = entry.Type
					if details, ok := entry.Strategy.Extract(ctx, id, rawContainer, base, f.client).(DetailProvider); ok {
						specific = details
					}
					break
				}
			}
			ch <- result{info: ContainerInfo{Type: matchedType, BaseContainerInfo: BaseContainerInfo(base), Specific: specific}, err: nil}
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
	if f.cfg.SortByCPU {
		sort.Slice(out, func(i, j int) bool { return out[i].CPUPercent > out[j].CPUPercent })
	} else {
		sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Names[0]) < strings.ToLower(out[j].Names[0]) })
	}

	return out, nil
}

func (f *Fetcher) DomainContainers(ctx context.Context) ([]domain.Container, error) {
	legacy, err := f.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Container, 0, len(legacy))
	for _, c := range legacy {
		var details domain.DetailProvider
		if dp, ok := c.Specific.(DetailProvider); ok {
			details = dp
		}
		out = append(out, domain.Container{ID: c.ID, Names: c.Names, Image: c.Image, Status: c.Status, CPUPercent: c.CPUPercent, MemoryMB: c.Mem, Type: c.Type, Details: details})
	}
	return out, nil
}
