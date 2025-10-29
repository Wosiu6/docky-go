package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/Microsoft/go-winio"
)

type DockerClient struct {
	http *http.Client
	url  string
}

func NewClient() (*DockerClient, error) {
	var transport *http.Transport
	if runtime.GOOS == "windows" {
		pipePath := `\\.\\pipe\\docker_engine`
		dial := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return winio.DialPipeContext(ctx, pipePath)
		}
		transport = &http.Transport{DialContext: dial}
	} else {
		socketPath := "/var/run/docker.sock"
		dial := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		}
		transport = &http.Transport{DialContext: dial}
	}
	return &DockerClient{http: &http.Client{Transport: transport, Timeout: 5 * time.Second}, url: "http://docker"}, nil
}

func (c *DockerClient) Ping(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", c.url+"/_ping", nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return errors.New(strings.TrimSpace(string(b)))
	}
	return nil
}

func (c *DockerClient) ListContainers(ctx context.Context) ([]map[string]interface{}, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", c.url+"/containers/json?all=1", nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list error: %s", strings.TrimSpace(string(b)))
	}
	var out []map[string]interface{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *DockerClient) ContainerStats(ctx context.Context, id string, dest interface{}) error {
	url := fmt.Sprintf("%s/containers/%s/stats?stream=false", c.url, id)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("stats error: %s", strings.TrimSpace(string(b)))
	}
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(dest)
}

func (c *DockerClient) ContainerInspect(ctx context.Context, id string, dest interface{}) error {
	url := fmt.Sprintf("%s/containers/%s/json", c.url, id)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("inspect error: %s", strings.TrimSpace(string(b)))
	}
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(dest)
}

func (c *DockerClient) GetHttpClient() *http.Client {
	return c.http
}

func (c *DockerClient) GetUrl() string {
	return c.url
}
