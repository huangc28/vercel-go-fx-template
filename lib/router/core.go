package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Handler interface {
	RegisterRoute(r *chi.Mux)
	Handle(w http.ResponseWriter, r *http.Request)
}

func AsRoute(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(Handler)),
			fx.ResultTags(`group:"handlers"`),
		),
	)
}

type MuxParams struct {
	fx.In

	Logger   *zap.Logger
	Handlers []Handler `group:"handlers"`
}

func NewMux(params MuxParams) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	for _, h := range params.Handlers {
		h.RegisterRoute(r)
	}

	return r
}
