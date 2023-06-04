package middlewares

import (
	"net/http"

	"github.com/getsentry/sentry-go"
)

func PerformanceMonitoring() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// start performance monitoring
			span := sentry.StartSpan(r.Context(), "http.server", sentry.WithTransactionName(r.URL.Path))
			defer span.Finish()

			// set span in request context
			ctx := sentry.SetHubOnContext(r.Context(), sentry.CurrentHub().Clone())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
