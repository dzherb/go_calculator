package orchestrator

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				slog.Error(
					"Failed to close request body",
					slog.String("error", err.Error()),
				)
			}
		}(r.Body)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().
			Set(
				"Access-Control-Allow-Headers",
				"Content-Type,access-control-allow-origin, access-control-allow-headers",
			)
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func ensureMethodsMiddleware(
	methods ...string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, m := range methods {
				if r.Method == m {
					next.ServeHTTP(w, r)
					return
				}
			}

			w.WriteHeader(http.StatusMethodNotAllowed)
			writeError(
				w,
				fmt.Errorf(
					"expected one of the methods: %s",
					strings.Join(methods, ", "),
				),
			)
		})
	}
}
