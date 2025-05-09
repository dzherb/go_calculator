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
		Handler:      CommonMiddleware(mux),
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
	mux.Handle("/api/v1/auth/register",
		EnsureMethodsMiddleware(http.MethodPost)(
			http.HandlerFunc(RegisterHandler),
		),
	)
	mux.Handle("/api/v1/auth/login",
		EnsureMethodsMiddleware(http.MethodPost)(
			http.HandlerFunc(LoginHandler),
		),
	)
	mux.Handle("/api/v1/calculate",
		AuthRequired(
			EnsureMethodsMiddleware(http.MethodPost)(
				http.HandlerFunc(CalculateHandler),
			),
		),
	)
	mux.Handle("/api/v1/expressions",
		AuthRequired(
			EnsureMethodsMiddleware(http.MethodGet)(
				http.HandlerFunc(ExpressionsHandler),
			),
		),
	)
	mux.Handle("/api/v1/expressions/{id}",
		AuthRequired(
			EnsureMethodsMiddleware(http.MethodGet)(
				http.HandlerFunc(ExpressionHandler),
			),
		),
	)
}
