package docker

import "context"

type Service interface {
	Health(ctx context.Context) error
	Containers(ctx context.Context) ([]map[string]interface{}, error)
	Stats(ctx context.Context, id string, dest interface{}) error
	Inspect(ctx context.Context, id string, dest interface{}) error
}

type serviceImpl struct { client DockerClient }

func NewService(c DockerClient) Service { return &serviceImpl{client: c} }

func (s *serviceImpl) Health(ctx context.Context) error { return s.client.Ping(ctx) }
func (s *serviceImpl) Containers(ctx context.Context) ([]map[string]interface{}, error) { return s.client.ListContainers(ctx) }
func (s *serviceImpl) Stats(ctx context.Context, id string, dest interface{}) error { return s.client.ContainerStats(ctx, id, dest) }
func (s *serviceImpl) Inspect(ctx context.Context, id string, dest interface{}) error { return s.client.ContainerInspect(ctx, id, dest) }
