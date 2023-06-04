package middlewares

import (
	"github.com/getsentry/sentry-go"
	"net/http"
)

func PerformanceMonitoring() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			println("PerformanceMonitoring")

			// print
			println("PerformanceMonitoring")

			// start performance monitoring
			span := sentry.StartSpan(r.Context(), "http.server", sentry.WithTransactionName(r.URL.Path))
			defer span.Finish()

			// set span in request context
			ctx := sentry.SetHubOnContext(r.Context(), sentry.CurrentHub().Clone())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
