package fx

import (
	"go.uber.org/fx"

	"github.com/example/vercel-go-service-template/cache"
	"github.com/example/vercel-go-service-template/config"
	"github.com/example/vercel-go-service-template/db"
	"github.com/example/vercel-go-service-template/lib/logs"
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
