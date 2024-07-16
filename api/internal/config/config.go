package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Port int    `env:"APP_PORT" env-required:"true"`
		Env  string `env:"APP_ENV" env-required:"true"`
	}

	DB struct {
		User string `env:"DB_USER" env-required:"true"`
		Pass string `env:"DB_PASS" env-required:"true"`
		Host string `env:"DB_HOST" env-required:"true"`
		Port int    `env:"DB_PORT" env-required:"true"`
		Name string `env:"DB_NAME" env-required:"true"`
	}

	PassportApi struct {
		Host string `env:"PASSPORT_API_HOST" env-required:"true"`
	}
}

func New() *Config {
	config := &Config{}

	if err := cleanenv.ReadEnv(config); err != nil {
		header := "Environment variables of application"
		f := cleanenv.FUsage(os.Stdout, config, &header)
		f()
		panic(err)
	}

	slog.SetDefault(initLogger(config.App.Env))

	slog.Info("config parsed", slog.Any("cfg", config))

	return config
}

func initLogger(env string) *slog.Logger {
	switch strings.ToLower(env) {
	case "local":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case "prod":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	}
}
