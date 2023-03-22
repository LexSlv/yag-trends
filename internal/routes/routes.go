package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/LexSlv/yagTrends/internal/handler"
)

func NewRoutes(h *handler.Handler) *fiber.App {
	app := fiber.New()
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		log.Info().Msg("Healthcheck is working!")
		return c.SendString("OK")
	})
	app.Get("/trends", h.GetGamesTopList)
	return app
}
