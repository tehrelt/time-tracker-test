package app

import (
	"em-test/internal/adapters"
	"em-test/internal/config"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	cfg  *config.Config
	http *fiber.App

	uc *adapters.UsersAdapter
	ac *adapters.ActivityAdapter
}

func New(cfg *config.Config, user *adapters.UsersAdapter, activity *adapters.ActivityAdapter) *App {

	http := fiber.New(fiber.Config{
		CaseSensitive: false,
	})

	return &App{
		cfg:  cfg,
		http: http,
		uc:   user,
		ac:   activity,
	}
}

func (a *App) initRoutes() {
	v1 := a.http.Group("/api/v1")

	users := v1.Group("/users")
	users.Get("/", a.uc.GetUsers())
	users.Post("/", a.uc.AddUser())

	activities := v1.Group("/activities")
	activities.Post("/", a.ac.Start())
	activities.Patch("/", a.ac.Stop())
	activities.Get("/:user_id", a.ac.GetSummary())
}

func (a *App) Run() error {
	a.initRoutes()
	return a.http.Listen(fmt.Sprintf(":%d", a.cfg.App.Port))
}
