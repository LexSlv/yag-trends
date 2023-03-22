package main

import (
	"context"
	"go.uber.org/fx"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/LexSlv/yagTrends/internal/config"
	"github.com/LexSlv/yagTrends/internal/cron"
	"github.com/LexSlv/yagTrends/internal/handler"
	"github.com/LexSlv/yagTrends/internal/repository"
	"github.com/LexSlv/yagTrends/internal/routes"
)

const serverPort = "8088"

func main() {
	app := fx.New(
		fx.Provide(
			config.NewConfig,
			repository.NewRepository,
			handler.NewHandler,
			routes.NewRoutes,
			cron.NewCron,
		),
		fx.Invoke(
			cron.RunCron,
			customizeLogger,
			setupLifeCycle),
	)

	app.Run()
	if err := app.Err(); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func setupLifeCycle(lc fx.Lifecycle, app *fiber.App) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Info().Msg("Server started on :" + serverPort)
				err := app.Listen(":" + serverPort)
				if err != nil {
					log.Fatal().Msg("Server not stated!")
					return
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			return app.Shutdown()
		},
	})
}

func customizeLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
