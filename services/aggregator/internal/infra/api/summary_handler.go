package api

import (
	"aggregator/internal/domain/entities"
	"aggregator/internal/infra/api/responses"
	"aggregator/internal/usecase/apperrors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type SummaryGetter interface {
	Handle(ctx context.Context, developerID string) (entities.DeveloperSummary, error)
}

type SummaryHandler struct {
	summaryGetter SummaryGetter
}

func NewSummaryHandler(summaryGetter SummaryGetter) *SummaryHandler {
	return &SummaryHandler{summaryGetter: summaryGetter}
}

func (h *SummaryHandler) Handle(c *fiber.Ctx) error {
	developerID := c.Params("developer_id")
	summary, err := h.summaryGetter.Handle(c.UserContext(), developerID)
	if errors.Is(err, apperrors.ErrNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "summary not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewDeveloperSummary(summary))
}
