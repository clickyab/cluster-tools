package errhandler

import (
	"context"
	"net/http"
)

// this route is used to health check this app only, not the database and other think.
func healthz(_ context.Context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
