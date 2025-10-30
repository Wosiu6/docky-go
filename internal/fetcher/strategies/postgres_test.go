package strategies

import (
	"context"
	"testing"

	"github.com/wosiu6/docky-go/internal/model"
)

type mockInspectClient struct{}

func (m *mockInspectClient) ContainerInspect(ctx context.Context, id string, v interface{}) error {
	out := v.(*struct {
		Config struct {
			Env []string `json:"Env"`
		} `json:"Config"`
		NetworkSettings struct {
			Ports map[string][]struct {
				HostPort string `json:"HostPort"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	})
	out.Config.Env = []string{
		"POSTGRES_DB=mydb",
		"POSTGRES_USER=admin",
		"POSTGRES_SSL_MODE=disable",
		"PGDATA=/var/lib/postgresql/data",
		"POSTGRES_MAX_CONNECTIONS=200",
	}
	out.NetworkSettings.Ports = map[string][]struct {
		HostPort string `json:"HostPort"`
	}{
		"5432/tcp": {{HostPort: "5432"}},
	}
	return nil
}

func TestPostgreSqlStrategy_Extract(t *testing.T) {
	s := &PostgreSqlStrategy{}
	base := model.BaseContainerInfo{ID: "x", Image: "postgres", Names: []string{"/pg"}}
	raw := map[string]interface{}{}
	res := s.Extract(context.Background(), "x", raw, base, &mockInspectClient{})
	info, ok := res.(*PostgreSqlContainerInfo)
	if !ok {
		t.Fatalf("expected PostgreSqlContainerInfo, got %T", res)
	}
	if info.Database != "mydb" {
		t.Errorf("Database parse failed: %s", info.Database)
	}
	if info.User != "admin" {
		t.Errorf("User parse failed: %s", info.User)
	}
	if info.SSLMode != "disable" {
		t.Errorf("SSLMode parse failed: %s", info.SSLMode)
	}
	if info.PGData != "/var/lib/postgresql/data" {
		t.Errorf("PGDATA parse failed: %s", info.PGData)
	}
	if info.MaxConnections != 200 {
		t.Errorf("MaxConnections parse failed: %d", info.MaxConnections)
	}
	if info.Port != 5432 {
		t.Errorf("Port parse failed: %d", info.Port)
	}
	fields := info.DetailFields()
	if fields["Database"] != "mydb" || fields["User"] != "admin" || fields["Port"] != "5432" {
		t.Errorf("unexpected detail fields: %#v", fields)
	}
}
