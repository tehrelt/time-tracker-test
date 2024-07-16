package adapters

import (
	"em-test/internal/domain"
	"em-test/internal/lib/filters"
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ActivityService interface {
	Start(userId string) error
	Stop(userId string) error
	GetSummary(*filters.Activity) (*domain.ActivitySummary, error)
}

type ActivityAdapter struct {
	activityService ActivityService
}

func NewActivityAdapter(activityService ActivityService) *ActivityAdapter {
	return &ActivityAdapter{
		activityService: activityService,
	}
}

func (a *ActivityAdapter) Start() fiber.Handler {

	type request struct {
		UserId string `json:"userId"`
	}

	return func(c *fiber.Ctx) error {

		req := new(request)
		if err := c.BodyParser(req); err != nil {
			return internal(c, fiber.Map{
				"error": err.Error(),
			})
		}

		if err := a.activityService.Start(req.UserId); err != nil {
			if errors.Is(err, domain.ErrUserAlreadyWorking) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			return internal(c, fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "activity started",
		})
	}
}

func (a *ActivityAdapter) Stop() fiber.Handler {
	type request struct {
		UserId string `json:"userId"`
	}

	return func(c *fiber.Ctx) error {

		req := new(request)
		if err := c.BodyParser(req); err != nil {
			return internal(c, fiber.Map{
				"error": err.Error(),
			})
		}

		if err := a.activityService.Stop(req.UserId); err != nil {
			if errors.Is(err, domain.ErrUserNotWorking) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return internal(c, fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "activity finished",
		})
	}

}

func (a *ActivityAdapter) GetSummary() fiber.Handler {

	fn := "ActivityAdapter.GetSummary"
	logger := slog.With(slog.String("fn", fn))

	return func(c *fiber.Ctx) error {

		userId := c.Params("user_id")
		startTime := c.Query("start_time")
		endTime := c.Query("end_time")

		filters := &filters.Activity{
			UserId: userId,
		}

		if startTime != "" {
			time, err := time.Parse("dd:MM:YYYY-HH:MM", startTime)
			if err != nil {
				logger.Error("cannot parse startTime", slog.String("rawStartTime", startTime))
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			filters.StartTime = &time
		}

		if endTime != "" {
			time, err := time.Parse("dd:MM:YYYY-HH:MM", endTime)
			if err != nil {
				logger.Error("cannot parse endTime", slog.String("rawEndTime", endTime))
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			filters.EndTime = &time
		}

		logger.Debug("filters setup", slog.Any("filters", filters))

		summary, err := a.activityService.GetSummary(filters)
		if err != nil {
			logger.Error("failed to get summary", slog.Any("filters", filters), slog.String("err", err.Error()))
			return internal(c, fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(summary)
	}
}
