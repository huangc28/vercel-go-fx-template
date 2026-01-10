package db

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/example/vercel-go-service-template/config"
)

func NewSQLXPostgresDB(lc fx.Lifecycle, cfg config.Config, logger *zap.Logger) (*sqlx.DB, error) {
	if cfg.PGURL == "" {
		logger.Info("postgres disabled (pg_url not set)")
		return nil, nil
	}

	db, err := sqlx.Open("pgx", cfg.PGURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}
