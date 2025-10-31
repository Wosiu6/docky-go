package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DockerClient interface {
	Ping(ctx context.Context) error
	ListContainers(ctx context.Context) ([]map[string]any, error)
	ContainerStats(ctx context.Context, id string, dest any) error
	ContainerInspect(ctx context.Context, id string, dest any) error
	GetHttpClient() *http.Client
	GetUrl() string
}

type dockerClientImpl struct {
	http *http.Client
	url  string
}

type Option func(*clientOptions)

type clientOptions struct {
	timeout time.Duration
	baseURL string
}

func WithTimeout(d time.Duration) Option { return func(o *clientOptions) { o.timeout = d } }
func WithBaseURL(u string) Option        { return func(o *clientOptions) { o.baseURL = u } }

func NewClient() (DockerClient, error) { return NewClientWithOptions() }
func NewClientWithOptions(opts ...Option) (DockerClient, error) {
	cfg := clientOptions{timeout: 5 * time.Second, baseURL: "http://docker"}
	for _, opt := range opts {
		opt(&cfg)
	}
	transport := buildTransport()
	return &dockerClientImpl{http: &http.Client{Transport: transport, Timeout: cfg.timeout}, url: cfg.baseURL}, nil
}

type HTTPError struct {
	Op     string
	Status int
	Body   string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("docker %s failed: status=%d body=%s", e.Op, e.Status, e.Body)
}

func (c *dockerClientImpl) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+"/_ping", nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return &HTTPError{Op: "ping", Status: resp.StatusCode, Body: strings.TrimSpace(string(b))}
	}
	return nil
}

func (c *dockerClientImpl) ListContainers(ctx context.Context) ([]map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+"/containers/json?all=1", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return nil, &HTTPError{Op: "list", Status: resp.StatusCode, Body: strings.TrimSpace(string(b))}
	}
	var out []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dockerClientImpl) ContainerStats(ctx context.Context, id string, dest any) error {
	url := fmt.Sprintf("%s/containers/%s/stats?stream=false", c.url, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return &HTTPError{Op: "stats", Status: resp.StatusCode, Body: strings.TrimSpace(string(b))}
	}
	return json.NewDecoder(resp.Body).Decode(dest)
}

func (c *dockerClientImpl) ContainerInspect(ctx context.Context, id string, dest any) error {
	url := fmt.Sprintf("%s/containers/%s/json", c.url, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return &HTTPError{Op: "inspect", Status: resp.StatusCode, Body: strings.TrimSpace(string(b))}
	}
	return json.NewDecoder(resp.Body).Decode(dest)
}

func (c *dockerClientImpl) GetHttpClient() *http.Client { return c.http }
func (c *dockerClientImpl) GetUrl() string              { return c.url }
