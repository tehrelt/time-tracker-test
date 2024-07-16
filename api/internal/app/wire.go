//go:build wireinject
// +build wireinject

package app

import (
	"em-test/internal/adapters"
	"em-test/internal/config"
	"em-test/internal/repositories"
	"em-test/internal/services"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

func Init() (*App, func(), error) {
	panic(wire.Build(
		New,
		wire.NewSet(config.New),
		wire.NewSet(initDB),

		wire.NewSet(repositories.NewUsersRepository),
		wire.NewSet(repositories.NewActivityRepository),
		wire.NewSet(repositories.NewPassportApi),

		wire.Bind(new(services.UserRepository), new(*repositories.UsersRepository)),
		wire.Bind(new(services.UserFinder), new(*repositories.PassportApi)),
		wire.Bind(new(services.ActivityRepository), new(*repositories.ActivityRepository)),

		wire.NewSet(services.NewUserService),
		wire.NewSet(services.NewActivityService),

		wire.Bind(new(adapters.UsersService), new(*services.UsersService)),
		wire.Bind(new(adapters.ActivityService), new(*services.ActivityService)),

		wire.NewSet(adapters.NewUsersAdapter),
		wire.NewSet(adapters.NewActivityAdapter),
	))
}

func initDB(cfg *config.Config) (*sqlx.DB, func(), error) {

	host := cfg.DB.Host
	port := cfg.DB.Port
	user := cfg.DB.User
	pass := cfg.DB.Pass
	name := cfg.DB.Name

	cs := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, pass, host, port, name)

	log.Printf("connecting to %s\n", cs)

	db, err := sqlx.Open("postgres", cs)
	if err != nil {
		return nil, nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, func() {
			db.Close()
		}, err
	}

	return db, func() { db.Close() }, nil
}
