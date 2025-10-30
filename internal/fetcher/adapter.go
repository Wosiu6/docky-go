package fetcher

import (
	"context"

	"github.com/wosiu6/docky-go/internal/domain"
)

type ServiceAdapter struct{ f *Fetcher }

func NewServiceAdapter(f *Fetcher) *ServiceAdapter { return &ServiceAdapter{f: f} }
func (a *ServiceAdapter) FetchAll(ctx context.Context) ([]domain.Container, error) {
	return a.f.DomainContainers(ctx)
}
