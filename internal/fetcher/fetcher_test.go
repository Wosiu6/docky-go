package fetcher

import (
	"context"
	"net/http"
	"testing"
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

func TestFetcher_FetchAll(t *testing.T) {
	f := New(&mockDockerClient{})
	ctx := context.Background()
	items, err := f.FetchAll(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
	if items[0].Type != TypePostgreSQL && items[0].Type != TypeGeneric {
		t.Errorf("unexpected container type: %s", items[0].Type)
	}
}
