package main

import (
	"flag"
	"github.com/pedromspeixoto/posts-api/internal/pkg/sentry"

	"github.com/pedromspeixoto/posts-api/internal/config"
	"github.com/pedromspeixoto/posts-api/internal/data"
	"github.com/pedromspeixoto/posts-api/internal/data/models"
	"github.com/pedromspeixoto/posts-api/internal/domain"
	"github.com/pedromspeixoto/posts-api/internal/http"
	"github.com/pedromspeixoto/posts-api/internal/http/handlers"
	"github.com/pedromspeixoto/posts-api/internal/pkg/logger"
	"github.com/pedromspeixoto/posts-api/internal/pkg/validator"
	"go.uber.org/fx"
)

// @title Posts API
// @version 1.0
// @description Posts API - Create blog posts and store in database
// @BasePath /
func main() {
	var cfgFilePath string
	flag.StringVar(
		&cfgFilePath,
		"config",
		"",
		"Path to config file. If not provided, config will be parsed from the environment.",
	)
	flag.Parse()

	app := fx.New(
		// Provide
		config.ProvideConfig(cfgFilePath),
		logger.ProvideLogger(),
		validator.ProvideValidator(),
		sentry.ProvideSentry(),
		data.ProvideData(),
		models.ProvideModels(),
		domain.ProvideDomains(),
		handlers.ProvideHandlers(),
		// Invoke
		http.InvokeServer(),
	)

	app.Run()
}
