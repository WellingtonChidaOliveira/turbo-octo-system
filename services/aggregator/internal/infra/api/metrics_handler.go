package api

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/infra/api/responses"
	"context"

	"github.com/gofiber/fiber/v2"
)

type EventsByDeveloperGetter interface {
	Handle(ctx context.Context, developerID string) ([]entities.ProcessedEvent, error)
}

type MetricsHandler struct {
	eventsGetter EventsByDeveloperGetter
}

func NewMetricsHandler(eventsGetter EventsByDeveloperGetter) *MetricsHandler {
	return &MetricsHandler{eventsGetter: eventsGetter}
}

func (h *MetricsHandler) Handle(c *fiber.Ctx) error {
	developerID := c.Params("developer_id")
	events, err := h.eventsGetter.Handle(c.UserContext(), developerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewProcessedEvents(events))
}
