package strategies

import "testing"

func TestRegistry_NotEmpty(t *testing.T) {
	entries := Registry()
	if len(entries) == 0 {
		t.Fatal("registry should not be empty")
	}
	foundPostgres := false
	foundRedis := false
	for _, e := range entries {
		if e.Type == "postgresql" {
			foundPostgres = true
		}
		if e.Type == "redis" {
			foundRedis = true
		}
	}
	if !foundPostgres {
		t.Error("postgres strategy missing")
	}
	if !foundRedis {
		t.Error("redis strategy missing")
	}
}

func TestPostgresStrategy_Match(t *testing.T) {
	s := &PostgreSqlStrategy{}
	cases := []struct {
		image string
		want  bool
	}{
		{"postgres", true},
		{"library/postgres:14", true},
		{"my-postgresql-custom", true},
		{"redis", false},
		{"nginx", false},
	}
	for _, c := range cases {
		if got := s.Match(c.image); got != c.want {
			t.Errorf("Match(%q) = %v want %v", c.image, got, c.want)
		}
	}
}

func TestRegistryUniqueness(t *testing.T) {
	entries := Registry()
	if len(entries) == 0 {
		t.Fatalf("registry should not be empty")
	}
	seen := map[string]struct{}{}
	for _, e := range entries {
		if _, ok := seen[string(e.Type)]; ok {
			t.Fatalf("duplicate container type in registry: %s", e.Type)
		}
		seen[string(e.Type)] = struct{}{}
		if e.Strategy == nil {
			t.Fatalf("nil strategy for type %s", e.Type)
		}
	}
}
