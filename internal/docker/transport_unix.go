//go:build !windows

package docker

import (
	"context"
	"net"
	"net/http"
)

func buildTransport() *http.Transport {
	socketPath := "/var/run/docker.sock"
	dial := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return net.Dial("unix", socketPath)
	}
	return &http.Transport{DialContext: dial}
}
