package postgres

import (
	"blog/internal/config"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
)

func NewPostgresDB(cfg *config.Config) *sqlx.DB {
	storagePath := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.DbName,
		cfg.Postgres.Password,
	)

	db, err := sqlx.Connect("pgx", storagePath)
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("cannot ping database: %s", err)
	}

	return db
}
