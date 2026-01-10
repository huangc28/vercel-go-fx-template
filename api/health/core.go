package health

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/fx"

	appfx "github.com/example/vercel-go-service-template/lib/app/fx"
	healthapp "github.com/example/vercel-go-service-template/lib/app/health"
	"github.com/example/vercel-go-service-template/lib/router"
	routerfx "github.com/example/vercel-go-service-template/lib/router/fx"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var mux http.Handler

	app := fx.New(
		appfx.CoreAppOptions,
		routerfx.CoreRouterOptions,
		router.AsRoute(healthapp.NewHandler),
		fx.Populate(&mux),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = app.Stop(stopCtx)
	}()

	mux.ServeHTTP(w, r)
}
