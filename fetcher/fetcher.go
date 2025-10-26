package fetcher

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/wosiu6/docky-go/docker"
)

type ContainerInfo struct {
	ID         string
	Names      []string
	Image      string
	CPUPercent float64
	Mem        uint64
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
		names := make([]string, 0, len(namesIface))
		for _, ni := range namesIface {
			if s, ok := ni.(string); ok {
				names = append(names, s)
			}
		}

		wg.Add(1)
		go func(id string, names []string, image string) {
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
				ch <- result{info: ContainerInfo{ID: id, Names: names, Image: image}, err: nil}
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

			ch <- result{info: ContainerInfo{ID: id, Names: names, Image: image, CPUPercent: cpu, Mem: v.MemoryStats.Usage}, err: nil}
		}(id, names, image)
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
