package docker

import (
	"runtime"
	"testing"
)

func TestBuildTransport_NotNil(t *testing.T) {
	tr := buildTransport()
	if tr == nil {
		t.Fatal("expected non-nil transport")
	}
	// DialContext should be set
	if tr.DialContext == nil {
		t.Fatal("expected DialContext to be configured")
	}
	// basic OS assertion (no panic)
	_ = runtime.GOOS
}
