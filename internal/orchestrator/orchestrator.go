package orchestrator

import (
	"context"
	"time"

	"github.com/wosiu6/docky-go/internal/domain"
	ilog "github.com/wosiu6/docky-go/internal/log"
)

type FetchService interface {
	FetchAll(ctx context.Context) ([]domain.Container, error)
}

type UiApp interface {
	SetData([]domain.Container)
	Run() error
}

type Orchestrator struct {
	fetch    FetchService
	ui       UiApp
	logger   ilog.Logger
	interval time.Duration
}

func New(fetch FetchService, ui UiApp, logger ilog.Logger, interval time.Duration) *Orchestrator {
	return &Orchestrator{fetch: fetch, ui: ui, logger: logger, interval: interval}
}

func (o *Orchestrator) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() { errCh <- o.ui.Run() }()

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(300 * time.Millisecond):
			if err := o.refreshOnce(ctx); err != nil {
				o.logger.Error("initial fetch failed", "error", err)
			}
		}
	}()

	ticker := time.NewTicker(o.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			o.logger.Info("context canceled; stopping orchestrator")
			return ctx.Err()
		case <-ticker.C:
			if err := o.refreshOnce(ctx); err != nil {
				o.logger.Error("periodic fetch failed", "error", err)
			}
		case err := <-errCh:
			return err
		}
	}
}

func (o *Orchestrator) refreshOnce(ctx context.Context) error {
	containers, err := o.fetch.FetchAll(ctx)
	if err != nil {
		return err
	}

	o.ui.SetData(containers)
	return nil
}
