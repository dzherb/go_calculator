package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const readTimeout = 15 * time.Second
const writeTimeout = 15 * time.Second

func (a *App) ServeHTTP(ctx context.Context) error {
	mux := http.NewServeMux()
	registerHandlers(mux)

	addr := a.config.Host + ":" + a.config.Port

	srv := http.Server{
		Addr:         a.config.Host + ":" + a.config.Port,
		Handler:      commonMiddleware(mux),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	go func() {
		<-ctx.Done()

		err := srv.Shutdown(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	slog.Info(fmt.Sprintf("HTTP server is listening on %s", addr))

	err := srv.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(
			"HTTP server stopped with an error",
			"error", err,
		)

		return err
	}

	slog.Info("HTTP server stopped")

	return nil
}

func registerHandlers(mux *http.ServeMux) {
	mux.Handle("/api/v1/calculate",
		ensureMethodsMiddleware(http.MethodPost)(
			http.HandlerFunc(calculateHandler),
		),
	)
	mux.Handle("/api/v1/expressions",
		ensureMethodsMiddleware(http.MethodGet)(
			http.HandlerFunc(expressionsHandler),
		),
	)
	mux.Handle("/api/v1/expressions/{id}",
		ensureMethodsMiddleware(http.MethodGet)(
			http.HandlerFunc(expressionHandler),
		),
	)
}
