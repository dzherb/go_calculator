package orchestrator

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dzherb/go_calculator/calculator/pkg/security"
)

func CommonMiddleware(next http.Handler) http.Handler {
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

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("recovered from panic", "error", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func EnsureMethodsMiddleware(
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
			WriteError(
				w,
				fmt.Errorf(
					"expected one of the methods: %s",
					strings.Join(methods, ", "),
				),
			)
		})
	}
}

type ctxKey string

const UserIDKey ctxKey = "userID"
const TokenPrefix = "Bearer "

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, TokenPrefix) {
			w.WriteHeader(http.StatusUnauthorized)
			WriteError(w, fmt.Errorf("token must be provided"))

			return
		}

		token := strings.TrimPrefix(authHeader, TokenPrefix)

		userID, err := security.ValidateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			WriteError(w, err)
		}

		r = r.WithContext(context.WithValue(r.Context(), UserIDKey, userID))
		next.ServeHTTP(w, r)
	})
}
