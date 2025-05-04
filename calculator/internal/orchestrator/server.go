package orchestrator

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const readTimeout = 15 * time.Second
const writeTimeout = 15 * time.Second

func (a *Application) RunServer() error {
	http.Handle("/api/v1/calculate", commonMiddleware(
		ensureMethodsMiddleware(http.MethodPost)(
			http.HandlerFunc(calculateHandler),
		),
	),
	)
	http.Handle("/api/v1/expressions", commonMiddleware(
		ensureMethodsMiddleware(http.MethodGet)(
			http.HandlerFunc(expressionsHandler),
		),
	),
	)
	http.Handle("/api/v1/expressions/{id}", commonMiddleware(
		ensureMethodsMiddleware(http.MethodGet)(
			http.HandlerFunc(expressionHandler),
		),
	),
	)

	http.Handle("/internal/task", commonMiddleware(
		ensureMethodsMiddleware(http.MethodGet, http.MethodPost)(
			http.HandlerFunc(taskHandler),
		),
	),
	)

	srv := http.Server{
		Addr:         a.config.Addr + ":" + a.config.Port,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	slog.Info(fmt.Sprintf("Listening on port %s", a.config.Port))

	return srv.ListenAndServe()
}
