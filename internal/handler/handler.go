package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/LexSlv/yagTrends/internal/config"
	"github.com/LexSlv/yagTrends/internal/models"
	"github.com/LexSlv/yagTrends/internal/repository"
)

type Handler struct {
	pool *repository.Repo
	cnf  *config.Cnf
}

func NewHandler(pool *repository.Repo, cnf *config.Cnf) *Handler {
	return &Handler{
		pool: pool,
		cnf:  cnf,
	}
}

func (h *Handler) GetGamesTopList(c *fiber.Ctx) error {
	offset := "0"
	limit := "100"
	if c.Query("offset") != "" && c.Query("limit") != "" {
		offset = c.Query("offset")
		limit = c.Query("limit")
	}

	games, err := h.pool.GetGamesTrends(offset, limit)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	return c.JSON(games)
}

func (h *Handler) CollectDevelopers() {
	log.Info().Msg("[CRON] CollectDevelopers: START")
	resp, err := http.Get(h.cnf.Links.MainFeedLink)
	if err != nil {
		log.Error().Msg("Http resp error")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Error().Msg("Body close error")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	var feedData models.FeedDevelopers
	err = json.Unmarshal(body, &feedData)
	if err != nil {
		log.Error().Msg("Feed unmarshal error")
	}

	var wg sync.WaitGroup
	developers := make(map[string]interface{})
	for _, developer := range feedData.Feed {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, item := range developer.Items {
				developers[item.Developer.Name] = nil
			}
		}()
		wg.Wait()
	}

	err = h.pool.CollectDeveloper(developers)
	if err != nil {
		log.Error().Msg("CollectDeveloper: " + err.Error())
	}

	log.Info().Msg("[CRON] CollectDevelopers: END")
}

func (h *Handler) CollectDeveloperTrends() {
	log.Info().Msg("[CRON] CollectDeveloperTrends: START")
	developer, err := h.pool.GetRandomDeveloper()
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	developer = url.QueryEscape(developer)

	parseUrl := h.cnf.Links.DeveloperFeedLink + developer
	resp, err := http.Get(parseUrl)
	if err != nil {
		log.Error().Msg("Http resp error")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Error().Msg("Body close error")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	var feedData models.FeedGames
	err = json.Unmarshal(body, &feedData)
	if err != nil {
		log.Error().Msg("Feed unmarshal error")
	}

	err = h.pool.CollectGames(feedData)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	log.Info().Msg("[CRON] CollectDeveloperTrends: END")
}
