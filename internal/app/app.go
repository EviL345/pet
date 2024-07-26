package app

import (
	"blog/internal/config"
	"github.com/jmoiron/sqlx"
)

type App struct {
	cfg *config.Config
	db  *sqlx.DB
}

func New(cfg *config.Config, db *sqlx.DB) *App {
	return &App{cfg: cfg, db: db}
}

func (a *App) Run() error {
	if err := a.mapHandlers(); err != nil {
		return err
	}

	return nil
}

func (a *App) mapHandlers() error {

	return nil
}
