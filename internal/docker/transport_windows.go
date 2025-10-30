//go:build windows

package docker

import (
	"context"
	"net"
	"net/http"

	"github.com/Microsoft/go-winio"
)

func buildTransport() *http.Transport {
	pipePath := `\\.\\pipe\\docker_engine`
	dial := func(ctx context.Context, network, addr string) (net.Conn, error) {
		ch := make(chan struct {
			c   net.Conn
			err error
		}, 1)
		go func() {
			c, err := winio.DialPipe(pipePath, nil)
			ch <- struct {
				c   net.Conn
				err error
			}{c, err}
		}()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case res := <-ch:
			return res.c, res.err
		}
	}
	return &http.Transport{DialContext: dial}
}
