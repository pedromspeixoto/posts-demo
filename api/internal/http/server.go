package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/pedromspeixoto/posts-api/internal/pkg/sentry"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/pedromspeixoto/posts-api/docs"
	"github.com/pedromspeixoto/posts-api/internal/config"
	"github.com/pedromspeixoto/posts-api/internal/http/handlers/health"
	"github.com/pedromspeixoto/posts-api/internal/http/handlers/posts"
	"github.com/pedromspeixoto/posts-api/internal/http/middlewares"
	"github.com/pedromspeixoto/posts-api/internal/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/fx"
)

func InvokeServer() fx.Option {
	return fx.Invoke(NewHTTPServer)
}

type serverDependencies struct {
	fx.In

	LifeCycle            fx.Lifecycle
	Config               *config.Config
	Logger               *logger.LoggingClient
	Sentry               *sentry.Sentry
	HealthServiceHandler health.HealthServiceHandler
	PostServiceHandler   posts.PostServiceHandler
}

func NewHTTPServer(lc fx.Lifecycle, deps serverDependencies) *http.Server {

	// http server definition
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", deps.Config.Port),
	}

	// server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// listen for syscall signals for process to quit
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// set routes
	registerRoutes(server, deps)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			deps.Logger.GetLogger().Info("starting HTTP server")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			deps.Logger.GetLogger().Info("stopping HTTP server")
			go func() {
				<-signals

				// shutdown signal with grace period of 30 seconds
				shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)
				go func() {
					<-shutdownCtx.Done()
					if shutdownCtx.Err() == context.DeadlineExceeded {
						deps.Logger.GetLogger().Fatal("graceful shutdown timed out.. forcing exit.")
					}
				}()

				// trigger graceful shutdown
				err := server.Shutdown(shutdownCtx)
				if err != nil {
					deps.Logger.GetLogger().Fatal(err.Error())
				}
				serverStopCtx()
			}()
			<-serverCtx.Done()
			return nil
		},
	})

	return server
}

func registerRoutes(server *http.Server, deps serverDependencies) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middlewares.RequestsLogger(deps.Logger.GetLogger()))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// cors support
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "sentry-trace", "baggage"},
	}))

	// define Sentry middleware if Sentry is enabled
	if deps.Config.SentryEnabled {
		sentryMiddleware := sentryhttp.New(sentryhttp.Options{
			Repanic: true,
		})

		r.Use(sentryMiddleware.Handle)
	}

	// define performance middleware
	r.Use(middlewares.PerformanceMonitoring())

	// swagger
	r.Mount("/swagger", httpSwagger.WrapHandler)

	// health
	r.Mount("/health", deps.HealthServiceHandler.Routes())

	// panic
	r.Mount("/panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Hello from backend Sentry!")
	}))

	// routes
	r.Group(func(r chi.Router) {
		r.Mount("/v1/posts", deps.PostServiceHandler.Routes())
	})

	server.Handler = r
}
