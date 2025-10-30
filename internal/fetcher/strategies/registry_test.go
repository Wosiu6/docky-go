package strategies

import "testing"

func TestRegistryUniqueness(t *testing.T) {
    entries := Registry()
    if len(entries) == 0 { t.Fatalf("registry should not be empty") }
    seen := map[string]struct{}{}
    for _, e := range entries {
        if _, ok := seen[string(e.Type)]; ok {
            t.Fatalf("duplicate container type in registry: %s", e.Type)
        }
        seen[string(e.Type)] = struct{}{}
        if e.Strategy == nil { t.Fatalf("nil strategy for type %s", e.Type) }
    }
}
