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
	if tr.DialContext == nil {
		t.Fatal("expected DialContext to be configured")
	}
	_ = runtime.GOOS
}
