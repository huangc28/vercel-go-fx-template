package inngestclient

import (
	"context"
	"log/slog"

	"github.com/inngest/inngestgo"
	"go.uber.org/zap"

	"github.com/example/vercel-go-service-template/config"
)

// NewInngestClient creates an Inngest Go SDK client.
//
// Note: Event/signing keys are read from environment variables by default:
// - INNGEST_EVENT_KEY
// - INNGEST_SIGNING_KEY
// - INNGEST_SIGNING_KEY_FALLBACK (optional)
func NewInngestClient(cfg config.Config) (inngestgo.Client, error) {
	return inngestgo.NewClient(inngestgo.ClientOpts{
		AppID:  cfg.InngestID,
		Logger: slog.Default(),
	})
}

// RegisterExampleCron registers a minimal example Inngest function.
func RegisterExampleCron(cli inngestgo.Client, logger *zap.Logger) error {
	_, err := inngestgo.CreateFunction(
		cli,
		inngestgo.FunctionOpts{
			ID:   "example-cron",
			Name: "Example Cron",
		},
		inngestgo.CronTrigger("0 * * * *"),
		func(ctx context.Context, _ inngestgo.Input[any]) (any, error) {
			logger.Info("running example cron")
			return map[string]any{"ok": true}, nil
		},
	)
	return err
}
