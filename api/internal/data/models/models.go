package models

import (
	"github.com/pedromspeixoto/posts-api/internal/data/models/posts"
	"go.uber.org/fx"
)

func ProvideModels() fx.Option {
	return fx.Options(
		fx.Provide(
			posts.NewPostRepository,
		),
	)
}
