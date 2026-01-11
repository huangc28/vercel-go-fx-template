package health

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/huangc28/vercel-go-fx-template/lib/pkg/render"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoute(r *chi.Mux) {
	r.Get("/", h.Handle)
	r.Get("/health", h.Handle)
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	render.ChiJSON(w, r, http.StatusOK, map[string]any{
		"ok": true,
	})
}
