package adapters

import (
	"em-test/internal/domain"
	"em-test/internal/lib/dto"
	"em-test/internal/lib/filters"
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UsersService interface {
	AddUser(dto *dto.AddUserDto) (*domain.User, error)
	GetUsers(filters *filters.UsersFilters) (users []*domain.User, total int64, err error)
}

type UsersAdapter struct {
	usersService UsersService
}

func NewUsersAdapter(usersService UsersService) *UsersAdapter {
	return &UsersAdapter{
		usersService: usersService,
	}
}

func (a *UsersAdapter) AddUser() fiber.Handler {
	type request struct {
		PassportInfo string `json:"passportNumber"`
	}

	return func(c *fiber.Ctx) error {
		var req request

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		parts := strings.Split(req.PassportInfo, " ")
		if len(parts) != 2 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "malformed passport info string",
			})
		}

		if len(parts[0]) != 4 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "unknown passport serie",
			})
		}

		if len(parts[1]) != 6 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "unknown passport number",
			})
		}

		serie, err := strconv.Atoi(parts[0])
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "unknown passport serie",
			})
		}

		number, err := strconv.Atoi(parts[1])
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "unknown passport number",
			})
		}

		user, err := a.usersService.AddUser(&dto.AddUserDto{
			PassportSerie:  serie,
			PassportNumber: number,
		})
		if err != nil {
			if errors.Is(err, domain.ErrUserAlreadyExists) {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "User already exists",
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user": user,
		})
	}
}

func (a *UsersAdapter) GetUsers() fiber.Handler {

	type response struct {
		Users []*domain.User `json:"users"`
		Total int64          `json:"count"`
	}

	fn := "UsersAdapter.GetUsers"
	logger := slog.With(slog.String("fn", fn))

	return func(c *fiber.Ctx) error {

		limit := c.QueryInt("limit")
		page := c.QueryInt("page", 1)
		offset := (page - 1) * limit
		surname := c.Query("surname")
		name := c.Query("name")
		address := c.Query("address")

		logger.Debug(
			"query params",
			slog.Int("limit", limit),
			slog.Int("page", page),
			slog.Int("offset", offset),
			slog.String("surname", surname),
			slog.String("name", name),
			slog.String("address", address),
		)

		filters := &filters.UsersFilters{
			Offset: &offset,
		}

		if limit != 0 {
			filters.Limit = &limit
		}

		if surname != "" {
			filters.Surname = &surname
		}
		if name != "" {
			filters.Name = &name
		}
		if address != "" {
			filters.Address = &address
		}

		users, total, err := a.usersService.GetUsers(filters)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(&response{
			Users: users,
			Total: total,
		})
	}
}
