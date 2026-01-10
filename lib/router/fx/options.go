package fx

import (
	"go.uber.org/fx"

	"github.com/example/vercel-go-service-template/lib/router"
)

var CoreRouterOptions = fx.Options(
	fx.Provide(router.NewMux),
)
