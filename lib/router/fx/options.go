package fx

import (
	"go.uber.org/fx"

	"github.com/huangc28/vercel-go-fx-template/lib/router"
)

var CoreRouterOptions = fx.Options(
	fx.Provide(router.NewMux),
)
