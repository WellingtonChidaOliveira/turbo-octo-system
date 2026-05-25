package api

import (
	"aggregator/internal/infra/api/responses"
	"aggregator/internal/usecase/ports"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct {
	queueHealth   ports.QueueHealthChecker
	storageHealth ports.StorageHealthChecker
}

func NewHealthHandler(queueHealth ports.QueueHealthChecker, storageHealth ports.StorageHealthChecker) *HealthHandler {
	return &HealthHandler{
		queueHealth:   queueHealth,
		storageHealth: storageHealth,
	}
}

func (h *HealthHandler) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()
	if err := h.queueHealth.CheckQueue(ctx); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(responses.Health{
			Status: "unhealthy",
			Error:  err.Error(),
		})
	}
	if err := h.storageHealth.CheckStorage(ctx); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(responses.Health{
			Status: "unhealthy",
			Error:  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.Health{Status: "ok"})
}
