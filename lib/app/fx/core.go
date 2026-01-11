package fx

import (
	"go.uber.org/fx"

	"github.com/huangc28/vercel-go-fx-template/cache"
	"github.com/huangc28/vercel-go-fx-template/config"
	"github.com/huangc28/vercel-go-fx-template/db"
	"github.com/huangc28/vercel-go-fx-template/lib/logs"
)

var CoreAppOptions = fx.Options(
	fx.Provide(
		config.NewViper,
		config.NewConfig,
		logs.NewLogger,
		db.NewSQLXPostgresDB,
		cache.NewRedis,
	),
)
