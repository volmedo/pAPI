package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
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
