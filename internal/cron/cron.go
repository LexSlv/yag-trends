package cron

import (
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"

	"github.com/LexSlv/yagTrends/internal/handler"
)

type Cron struct {
	cron *cron.Cron
}

func NewCron(h *handler.Handler) *Cron {
	scheduler, err := RunCron(h)
	if err != nil {
		log.Error().Msg("[CRON] NewCron: " + err.Error())
	}
	return &Cron{cron: scheduler}
}

func RunCron(h *handler.Handler) (*cron.Cron, error) {
	c := cron.New()
	_, err := c.AddFunc("0 */1 * * *", h.CollectDevelopers)
	if err != nil {
		log.Error().Msg("[CRON] CollectDevelopers: " + err.Error())
	}

	_, err = c.AddFunc("*/2 * * * *", h.CollectDeveloperTrends)
	if err != nil {
		log.Error().Msg("[CRON] CollectDeveloperTrends: " + err.Error())
	}
	c.Start()

	return c, nil
}
