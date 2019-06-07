package main

import (
	"net/http"

	"github.com/unrolled/recovery"
)

// newRecoveredHandler adds a basic panic recovery middleware so that clients
// get a 500 Internal Server Error when something goes wrong
func newRecoveredHandler(handler http.Handler) http.Handler {
	rec := recovery.New()
	return rec.Handler(handler)
}
