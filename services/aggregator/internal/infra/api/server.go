package api

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app     *fiber.App
	address string
}

func NewServer(
	address string,
	healthHandler *HealthHandler,
	metricsHandler *MetricsHandler,
	summaryHandler *SummaryHandler,
) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/health", healthHandler.Handle)
	app.Get("/metrics/:developer_id/summary", summaryHandler.Handle)
	app.Get("/metrics/:developer_id", metricsHandler.Handle)

	return &Server{
		app:     app,
		address: address,
	}
}

func (s *Server) Start() <-chan error {
	done := make(chan error, 1)
	go func() {
		done <- s.app.Listen(s.address)
	}()
	return done
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

func (s *Server) Test(req *http.Request) (*http.Response, error) {
	return s.app.Test(req)
}
