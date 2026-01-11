package health

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/fx"

	appfx "github.com/huangc28/vercel-go-fx-template/lib/app/fx"
	healthapp "github.com/huangc28/vercel-go-fx-template/lib/app/health"
	"github.com/huangc28/vercel-go-fx-template/lib/router"
	routerfx "github.com/huangc28/vercel-go-fx-template/lib/router/fx"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var mux http.Handler

	app := fx.New(
		appfx.CoreAppOptions,
		routerfx.CoreRouterOptions,
		router.AsRoute(healthapp.NewHandler),
		fx.Populate(&mux),
	)

	if err := app.Start(r.Context()); err != nil {
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
