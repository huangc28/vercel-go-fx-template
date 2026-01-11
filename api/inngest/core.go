package inngest

import (
	"context"
	"net/http"
	"time"

	"github.com/inngest/inngestgo"
	"go.uber.org/fx"

	appfx "github.com/huangc28/vercel-go-fx-template/lib/app/fx"
	"github.com/huangc28/vercel-go-fx-template/lib/pkg/inngestclient"
	"github.com/huangc28/vercel-go-fx-template/lib/pkg/render"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	app := fx.New(
		appfx.CoreAppOptions,
		fx.Provide(inngestclient.NewInngestClient),
		fx.Invoke(inngestclient.RegisterExampleCron),
		fx.Invoke(func(cli inngestgo.Client) {
			cli.Serve().ServeHTTP(w, r)
		}),
	)

	if err := app.Start(r.Context()); err != nil {
		render.ChiErr(w, r, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = app.Stop(stopCtx)
	}()
}
