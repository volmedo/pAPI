package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/unrolled/recovery"
)

// newMeasuredHandler creates a middleware that take essential metrics about
// the handler being measured, such as number of requests, duration of each request,
// concurrent or in-flight requests and response size.
// This function returns two handlers, the handler being measured and a Prometheus
// handler that exposes the metrics being collected
func newMeasuredHandler(handler http.Handler) (measuredH http.Handler, metricsH http.Handler) {
	recorder := metrics.NewRecorder(metrics.Config{
		Prefix: "pAPI",
	})
	mdlw := middleware.New(middleware.Config{
		Recorder: recorder,
	})

	return mdlw.Handler("", handler), promhttp.Handler()
}

// newRateLimitedHandler creates a new middleware based on ulule/limiter package that
// limits the request rate that is sent to the specified handler.
// The returned rate-limited handler will allow up to rps requests per second to
// handler. When the rate exceeds the limit, a "429 Too Many Requests" response will be
// sent back without invoking the wrapped handler.
func newRateLimitedHandler(rps int64, handler http.Handler) (http.Handler, error) {
	if rps <= 0 {
		return nil, fmt.Errorf("rps cannot be negative (rps = %d)", rps)
	}

	store := memory.NewStore()
	rate := limiter.Rate{
		Period: time.Second,
		Limit:  rps,
	}
	instance := limiter.New(store, rate)
	middleware := stdlib.NewMiddleware(instance)

	return middleware.Handler(handler), nil
}

// newRecoveredHandler adds a basic panic recovery middleware so that clients
// get a 500 Internal Server Error when something goes wrong
func newRecoverableHandler(handler http.Handler) http.Handler {
	rec := recovery.New()
	return rec.Handler(handler)
}

// newHealthHandler returns a basic health endpoint that can be used in readiness
// and liveness probes. It checks moving parts to report general availability
// of the service (currently, the connection with the DB is the only moving part).
// The health endpoint returns a 200 response when the service is available
// and a 500 one when it is not working correctly
func newHealthHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "cannot connect with DB: %v", err)
		} else {
			w.WriteHeader(200)
			fmt.Fprint(w, "ok")
		}
	})
}
