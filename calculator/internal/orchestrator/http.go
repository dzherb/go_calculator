package orchestrator

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const readTimeout = 15 * time.Second
const writeTimeout = 15 * time.Second

func (a *Application) ServeHTTP() error {
	mux := http.NewServeMux()
	registerHandlers(mux)

	addr := a.config.Host + ":" + a.config.Port

	srv := http.Server{
		Addr:         a.config.Host + ":" + a.config.Port,
		Handler:      commonMiddleware(mux),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	slog.Info(fmt.Sprintf("HTTP server is listening on %s", addr))

	return srv.ListenAndServe()
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
