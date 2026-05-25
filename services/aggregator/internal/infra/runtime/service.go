package runtime

import (
	"aggregator/internal/dto"
	"aggregator/internal/infra/api"
	"aggregator/internal/infra/config"
	"aggregator/internal/usecase/ports"
	"context"
	"log/slog"
	"time"
)

type Service struct {
	consumer ports.JobConsumer
	workers  ports.WorkerPool
	api      *api.Server
	cfg      config.Settings
}

func NewService(consumer ports.JobConsumer, workers ports.WorkerPool, apiServer *api.Server, cfg config.Settings) *Service {
	return &Service{
		consumer: consumer,
		workers:  workers,
		api:      apiServer,
		cfg:      cfg,
	}
}

func (s *Service) Run(ctx context.Context, stop context.CancelFunc) {
	jobs := make(chan dto.QueueMessage, s.cfg.Worker.JobBufferSize)
	workersDone := make(chan struct{})
	apiDone := s.api.Start()

	slog.Info("aggregator starting",
		"worker_count", s.cfg.Worker.Count,
		"job_buffer_size", s.cfg.Worker.JobBufferSize,
		"api_port", s.cfg.API.Port,
	)

	go s.consumer.Start(ctx, jobs)
	go func() {
		s.workers.Start(ctx, jobs)
		close(workersDone)
	}()

	select {
	case <-ctx.Done():
	case err := <-apiDone:
		if err != nil {
			slog.Error("api server stopped unexpectedly", "error", err)
		}
		stop()
	}

	s.shutdownAPI()
	s.awaitWorkers(workersDone)

	slog.Info("aggregator stopped")
}

func (s *Service) shutdownAPI() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.api.Shutdown(shutdownCtx); err != nil {
		slog.Error("api server shutdown failed", "error", err)
	}
}

func (s *Service) awaitWorkers(workersDone <-chan struct{}) {
	select {
	case <-workersDone:
	case <-time.After(s.cfg.ShutdownTimeout):
		slog.Warn("aggregator shutdown timeout exceeded",
			"shutdown_timeout", s.cfg.ShutdownTimeout.String(),
		)
	}
}
