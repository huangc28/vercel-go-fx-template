package inngest

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/fx"

	appfx "github.com/example/vercel-go-service-template/lib/app/fx"
	"github.com/example/vercel-go-service-template/lib/pkg/inngestclient"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler

	app := fx.New(
		appfx.CoreAppOptions,
		fx.Provide(inngestclient.NewHTTPHandler),
		fx.Populate(&handler),
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

	handler.ServeHTTP(w, r)
}
