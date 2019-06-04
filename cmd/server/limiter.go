package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

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
