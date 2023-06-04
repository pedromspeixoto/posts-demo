package domain

import (
	"go.uber.org/fx"

	"github.com/pedromspeixoto/posts-api/internal/domain/health"
	"github.com/pedromspeixoto/posts-api/internal/domain/posts"
)

func ProvideDomains() fx.Option {
	return fx.Provide(
		health.NewHealthService,
		posts.NewPostService,
	)
}
