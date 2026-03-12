package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/choegeun-won/terraform-gcp-wireguard-openclaw/apps/backend/internal/ingest"
)

type RetryWorker struct {
	service  *ingest.Service
	logger   *slog.Logger
	interval time.Duration
	limit    int
}

func NewRetryWorker(service *ingest.Service, logger *slog.Logger, interval time.Duration, limit int) *RetryWorker {
	if interval <= 0 {
		interval = time.Second
	}
	if limit <= 0 {
		limit = 10
	}

	return &RetryWorker{
		service:  service,
		logger:   logger,
		interval: interval,
		limit:    limit,
	}
}

func (w *RetryWorker) Run(ctx context.Context) {
	if w == nil || w.service == nil {
		return
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			result, err := w.service.ProcessDueRetries(ctx, w.limit)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				if w.logger != nil {
					w.logger.Warn("retry worker failed", "error", err)
				}
				continue
			}
			if w.logger != nil && result.Processed > 0 {
				w.logger.Info("retry worker processed ingest events",
					"processed", result.Processed,
					"completed", result.Completed,
					"rescheduled", result.Rescheduled,
					"dead_lettered", result.DeadLettered,
				)
			}
		}
	}
}
