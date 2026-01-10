package inngestclient

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/example/vercel-go-service-template/config"
	"github.com/example/vercel-go-service-template/lib/pkg/render"
)

// NewHTTPHandler is a minimal placeholder Inngest endpoint.
//
// Replace this with the official Inngest Go SDK integration if you want
// automatic function registration and event processing.
func NewHTTPHandler(cfg config.Config, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("inngest endpoint called")
		render.ChiJSON(w, r, http.StatusOK, map[string]any{
			"ok":        true,
			"app_id":    cfg.InngestID,
			"note":      "replace lib/pkg/inngestclient with official SDK integration",
			"path":      r.URL.Path,
			"method":    r.Method,
			"userAgent": r.UserAgent(),
		})
	})
}
