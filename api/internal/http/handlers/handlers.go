package handlers

import (
	"github.com/pedromspeixoto/posts-api/internal/http/handlers/health"
	"github.com/pedromspeixoto/posts-api/internal/http/handlers/posts"
	"go.uber.org/fx"
)

func ProvideHandlers() fx.Option {
	return fx.Provide(
		health.NewHealthServiceHandler,
		posts.NewPostServiceHandler,
	)
}
