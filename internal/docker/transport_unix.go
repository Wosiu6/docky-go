//go:build !windows

package docker

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// buildTransport attempts to discover a working Docker socket path.
// Order:
// 1. DOCKER_HOST=unix://...
// 2. /var/run/docker.sock
// 3. $XDG_RUNTIME_DIR/docker.sock
// 4. /run/user/$UID/docker.sock
// Returns transport that dials the first existing socket; if none exist we still
// return a transport whose dial will produce a descriptive error.
func buildTransport() *http.Transport {
	host := os.Getenv("DOCKER_HOST")
	var candidates []string
	if host != "" {
		if len(host) > 7 && host[:7] == "unix://" {
			candidates = append(candidates, host[7:])
		}
		// ignore tcp:// here; default http.Transport can handle it if we set baseURL accordingly elsewhere.
	}
	candidates = append(candidates, "/var/run/docker.sock")
	if xdg := os.Getenv("XDG_RUNTIME_DIR"); xdg != "" {
		candidates = append(candidates, filepath.Join(xdg, "docker.sock"))
	}
	if uid := os.Getuid(); uid > 0 {
		candidates = append(candidates, filepath.Join("/run/user", strconv.Itoa(uid), "docker.sock"))
	}

	chosen := ""
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && (fi.Mode()&os.ModeSocket) != 0 {
			chosen = c
			break
		}
	}
	if chosen == "" {
		// keep first candidate (maybe DOCKER_HOST) just for error path
		if len(candidates) > 0 {
			chosen = candidates[0]
		} else {
			chosen = "/var/run/docker.sock"
		}
	}
	if os.Getenv("DOCKY_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "[docky] using docker socket: %s\n", chosen)
	}
	dial := func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := net.Dial("unix", chosen)
		if err != nil {
			// augment error for clarity
			var e *net.OpError
			if errors.As(err, &e) {
				return nil, fmt.Errorf("docker socket dial failed (%s): %w", chosen, err)
			}
			return nil, fmt.Errorf("docker socket dial error (%s): %w", chosen, err)
		}
		return conn, nil
	}
	return &http.Transport{DialContext: dial}
}
