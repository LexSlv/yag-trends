package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/LexSlv/yagTrends/internal/config"
	"github.com/LexSlv/yagTrends/internal/models"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepository(cnf *config.Cnf) *Repo {
	pool, err := NewPgxPool(context.Background(), cnf)
	if err != nil {
		log.Error().Msg("[PGXPOOL]: " + err.Error())
	}

	return &Repo{db: pool}
}

func NewPgxPool(ctx context.Context, cnf *config.Cnf) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cnf.Db.User, cnf.Db.Pass, cnf.Db.Host, cnf.Db.Port, cnf.Db.Name)
	pgConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error().Msg("[PGXPOOL]: " + err.Error())
	}

	if err != nil {
		log.Error().Msg("[PGXPOOL]: " + err.Error())
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgConfig)
	if err != nil {
		log.Error().Msg("[PGXPOOL]: " + err.Error())
	}

	log.Info().Msg("Database connected!")

	return pool, nil
}

func (repo *Repo) CollectDeveloper(developers map[string]interface{}) error {
	for developer, _ := range developers {
		go func() {
			sqlStatement := `INSERT INTO developers (name) VALUES ($1) ON CONFLICT DO NOTHING`
			_, err := repo.db.Exec(context.Background(), sqlStatement, developer)
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}()
	}
	return nil
}

func (repo *Repo) GetRandomDeveloper() (string, error) {
	var developer string
	err := repo.db.QueryRow(context.Background(), "SELECT name FROM developers WHERE check_date < CURRENT_DATE OR check_date IS NULL ORDER BY random() LIMIT 1").Scan(&developer)
	if err != nil {
		log.Error().Msg("[PGXPOOL] GetRandomDeveloper select: " + err.Error())
		return "", err
	}

	_, err = repo.db.Exec(context.Background(), "UPDATE developers SET check_date=$1 WHERE name=$2", time.Now(), developer)
	if err != nil {
		log.Error().Msg("[PGXPOOL] GetRandomDeveloper update: " + err.Error())
		return "", err
	}

	return developer, nil
}

func (repo *Repo) CollectGames(games models.FeedGames) error {
	var wg sync.WaitGroup
	for _, game := range games.Feed {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, item := range game.Items {
				selectStmt := `SELECT players FROM games WHERE name = $1 AND developer = $2 ORDER BY created_at DESC LIMIT 1`
				playersOld := 0
				err := repo.db.QueryRow(context.Background(), selectStmt, item.Title, item.Developer.Name).Scan(&playersOld)
				if err != nil {
					log.Error().Msg("[PGXPOOL] CollectGames select: " + err.Error())
				}

				insertStmt := `INSERT INTO games (name, img, players, trend, developer, url) VALUES ($1, $2, $3, $4, $5, $6 ) RETURNING id`
				name := item.Title
				developer := item.Developer.Name
				url := item.URL
				img := item.Media.Icon.PrefixURL
				players := item.PlayersCount
				trend := item.PlayersCount - playersOld
				if playersOld == 0 {
					trend = 0
				}

				var id int
				err = repo.db.QueryRow(context.Background(), insertStmt, name, img, players, trend, developer, url).Scan(&id)
				if err != nil {
					log.Error().Msg("[PGXPOOL] CollectGames insert: " + err.Error())
				}
			}
		}()
		wg.Wait()
	}

	return nil
}

func (repo *Repo) GetGamesTrends(offset, limit string) ([]models.Game, error) {
	selectStmt := `SELECT id, name, img, players, trend, developer, created_at, url FROM games WHERE created_at = $1 AND trend > 0 ORDER BY trend DESC LIMIT $2 OFFSET $3`
	log.Info().Msg(offset)
	var games []models.Game
	rows, err := repo.db.Query(context.Background(), selectStmt, limit, offset)
	if err != nil {
		log.Error().Msg("[PGXPOOL] GetGamesTrends select: " + err.Error())
	}
	defer rows.Close()

	log.Info().Msg(selectStmt)

	for rows.Next() {
		var g models.Game
		err = rows.Scan(&g.ID, &g.Name, &g.Img, &g.Players, &g.Trend, &g.Developer, &g.CreatedAt, &g.URL)
		if err != nil {
			log.Error().Msg("[PGXPOOL] GetGamesTrends rows scan: " + err.Error())
		}
		games = append(games, g)
	}

	return games, err
}
