package fetcher

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/wosiu6/docky-go/internal/domain"
)

type mockDockerClient struct{}

func (m *mockDockerClient) Ping(ctx context.Context) error { return nil }
func (m *mockDockerClient) ListContainers(ctx context.Context) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{"Id": "abc123", "Names": []interface{}{"/test"}, "Image": "postgres", "State": "running", "Status": "Up"},
	}, nil
}
func (m *mockDockerClient) GetHttpClient() *http.Client { return nil }
func (m *mockDockerClient) GetUrl() string              { return "mock" }
func (m *mockDockerClient) ContainerInspect(ctx context.Context, id string, v interface{}) error {
	return nil
}
func (m *mockDockerClient) ContainerStats(ctx context.Context, id string, v interface{}) error {
	return nil
}

func TestFetcher_FetchAll_PostgresMatchFallback(t *testing.T) {
	f := New(&mockDockerClient{})
	ctx := context.Background()
	items, err := f.FetchAll(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Type != domain.ContainerTypePostgreSQL && items[0].Type != domain.ContainerTypeGeneric {
		t.Errorf("unexpected container type: %s", items[0].Type)
	}
	if len(items[0].Names) == 0 || items[0].Names[0] != "/test" {
		t.Errorf("expected name '/test', got %#v", items[0].Names)
	}
}

type mockDockerClientStats struct {
	statsCalls int
}

func (m *mockDockerClientStats) Ping(ctx context.Context) error { return nil }
func (m *mockDockerClientStats) ListContainers(ctx context.Context) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"Id": "id1", "Names": []interface{}{"/alpha"}, "Image": "redis", "State": "running"}}, nil
}
func (m *mockDockerClientStats) GetHttpClient() *http.Client { return nil }
func (m *mockDockerClientStats) GetUrl() string              { return "mock" }
func (m *mockDockerClientStats) ContainerInspect(ctx context.Context, id string, v interface{}) error {
	return nil
}
func (m *mockDockerClientStats) ContainerStats(ctx context.Context, id string, v interface{}) error {
	m.statsCalls++
	out := v.(*struct {
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
	})
	out.CPUStats.CPUUsage.TotalUsage = uint64(100 * m.statsCalls)
	out.CPUStats.CPUUsage.Percpu = []uint64{1, 2}
	out.CPUStats.SystemCPUUsage = uint64(1000 * m.statsCalls)
	out.CPUStats.OnlineCPUs = 2
	out.MemoryStats.Usage = 50 * 1024 * 1024 // 50MB
	return nil
}

func TestFetcher_CPUCalculation(t *testing.T) {
	m := &mockDockerClientStats{}
	f := New(m)
	ctx := context.Background()
	items, err := f.FetchAll(ctx)
	if err != nil {
		t.Fatalf("first fetch error: %v", err)
	}
	if items[0].CPUPercent != 0 {
		t.Errorf("expected 0 cpu on first snapshot, got %f", items[0].CPUPercent)
	}
	time.Sleep(10 * time.Millisecond)
	items2, err := f.FetchAll(ctx)
	if err != nil {
		t.Fatalf("second fetch error: %v", err)
	}
	if items2[0].CPUPercent <= 0 {
		t.Errorf("expected cpu > 0 on second snapshot, got %f", items2[0].CPUPercent)
	}
}

func TestServiceAdapter_DomainContainers(t *testing.T) {
	m := &mockDockerClient{}
	f := New(m)
	a := NewServiceAdapter(f)
	ctx := context.Background()
	containers, err := a.FetchAll(ctx)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(containers))
	}
	if containers[0].Type != domain.ContainerTypePostgreSQL && containers[0].Type != domain.ContainerTypeGeneric {
		t.Fatalf("unexpected type: %s", containers[0].Type)
	}
}

type mockDockerClientMulti struct{}

func (m *mockDockerClientMulti) Ping(ctx context.Context) error { return nil }
func (m *mockDockerClientMulti) ListContainers(ctx context.Context) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{"Id": "c1", "Names": []interface{}{"/b"}, "Image": "redis", "State": "running"},
		{"Id": "c2", "Names": []interface{}{"/A"}, "Image": "postgres", "State": "running"},
	}, nil
}
func (m *mockDockerClientMulti) GetHttpClient() *http.Client { return nil }
func (m *mockDockerClientMulti) GetUrl() string              { return "mock" }
func (m *mockDockerClientMulti) ContainerInspect(ctx context.Context, id string, v interface{}) error {
	return nil
}
func (m *mockDockerClientMulti) ContainerStats(ctx context.Context, id string, v interface{}) error {
	out := v.(*struct {
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
	})
	if id == "c1" {
		out.CPUStats.CPUUsage.TotalUsage = 100
	} else {
		out.CPUStats.CPUUsage.TotalUsage = 200
	}
	out.CPUStats.CPUUsage.Percpu = []uint64{50, 60}
	out.CPUStats.SystemCPUUsage = 1000
	out.CPUStats.OnlineCPUs = 2
	out.MemoryStats.Usage = 20 * 1024 * 1024
	return nil
}

func TestFetcher_SortByNameDefault(t *testing.T) {
	f := New(&mockDockerClientMulti{})
	items, err := f.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items")
	}
	if items[0].Names[0] != "/A" {
		t.Errorf("expected first name /A, got %s", items[0].Names[0])
	}
}

func TestFetcher_SortByCPU(t *testing.T) {
	cfg := FetcherConfig{Concurrency: 2, SortByCPU: true}
	f := NewWithConfig(&mockDockerClientMulti{}, cfg)
	_, _ = f.FetchAll(context.Background())
	time.Sleep(5 * time.Millisecond)
	items, err := f.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if items[0].CPUPercent < items[1].CPUPercent {
		t.Errorf("expected first item to have higher CPU")
	}
}

type mockDockerClientStatsError struct{}

func (m *mockDockerClientStatsError) Ping(ctx context.Context) error { return nil }
func (m *mockDockerClientStatsError) ListContainers(ctx context.Context) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"Id": "x", "Names": []interface{}{"/err"}, "Image": "unknown", "State": "running"}}, nil
}
func (m *mockDockerClientStatsError) GetHttpClient() *http.Client { return nil }
func (m *mockDockerClientStatsError) GetUrl() string              { return "mock" }
func (m *mockDockerClientStatsError) ContainerInspect(ctx context.Context, id string, v interface{}) error {
	return nil
}
func (m *mockDockerClientStatsError) ContainerStats(ctx context.Context, id string, v interface{}) error {
	return assertErr
}

var assertErr = &mockError{"stats failed"}

type mockError struct{ msg string }

func (e *mockError) Error() string { return e.msg }

func TestFetcher_StatsErrorFallback(t *testing.T) {
	f := New(&mockDockerClientStatsError{})
	items, err := f.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if items[0].Type != domain.ContainerTypeGeneric {
		t.Errorf("expected generic type on stats error")
	}
}

func TestFetcher_ConfigConcurrencyFallback(t *testing.T) {
	cfg := FetcherConfig{Concurrency: 0}
	f := NewWithConfig(&mockDockerClient{}, cfg)
	if f.cfg.Concurrency <= 0 {
		t.Errorf("expected concurrency fallback >0, got %d", f.cfg.Concurrency)
	}
}
