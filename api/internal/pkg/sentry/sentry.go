package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pedromspeixoto/posts-api/internal/config"
	"github.com/pedromspeixoto/posts-api/internal/pkg/logger"
	"go.uber.org/fx"
)

func ProvideSentry() fx.Option {
	return fx.Provide(
		NewSentryClient,
	)
}

type sentryDeps struct {
	fx.In

	Config *config.Config
	Logger *logger.LoggingClient
}

type Sentry struct {
	Deps   sentryDeps
	Logger logger.Logger
}

func NewSentryClient(deps sentryDeps) *Sentry {
	deps.Logger.GetLogger().Info("Initializing Sentry client")
	if !deps.Config.SentryEnabled {
		return nil
	}
	// check if sentry dsn is empty
	if deps.Config.SentryDSN == "" {
		deps.Logger.GetLogger().Warningf("sentry.Init: %s", "sentry dsn is empty")
		return nil
	}
	err := sentry.Init(sentry.ClientOptions{
		EnableTracing:    true,
		Dsn:              deps.Config.SentryDSN,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		deps.Logger.GetLogger().Warningf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	return &Sentry{
		Deps:   deps,
		Logger: deps.Logger.GetLogger(),
	}
}
