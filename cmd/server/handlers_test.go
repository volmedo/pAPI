package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

func TestHealth(t *testing.T) {
	goodConn, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock connection: %v", err)
	}

	badConn, err := sql.Open("postgres", "host=bad.host")
	if err != nil {
		t.Fatalf("Error creating bad connection: %v", err)
	}

	tests := map[string]struct {
		conn     *sql.DB
		wantCode int
	}{
		"healthy": {
			conn:     goodConn,
			wantCode: http.StatusOK,
		},
		"unhealthy": {
			conn:     badConn,
			wantCode: http.StatusInternalServerError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/health", nil)
			resp := httptest.NewRecorder()
			handler := newHealthHandler(tc.conn)
			handler.ServeHTTP(resp, req)

			if resp.Code != tc.wantCode {
				t.Fatalf("want %d but got %d", tc.wantCode, resp.Code)
			}
		})
	}
}
