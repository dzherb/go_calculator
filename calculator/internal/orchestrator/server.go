package orchestrator

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (a *Application) RunServer() error {
	http.Handle("/api/v1/calculate", commonMiddleware(
		ensureMethodsMiddleware("POST")(
			http.HandlerFunc(calculateHandler),
		),
	),
	)
	http.Handle("/api/v1/expressions", commonMiddleware(
		ensureMethodsMiddleware("GET")(
			http.HandlerFunc(expressionsHandler),
		),
	),
	)
	http.Handle("/api/v1/expressions/{id}", commonMiddleware(
		ensureMethodsMiddleware("GET")(
			http.HandlerFunc(expressionHandler),
		),
	),
	)

	http.Handle("/internal/task", commonMiddleware(
		ensureMethodsMiddleware("GET", "POST")(
			http.HandlerFunc(taskHandler),
		),
	),
	)

	slog.Info(fmt.Sprintf("Listening on port %s", a.config.Port))
	return http.ListenAndServe(a.config.Addr+":"+a.config.Port, nil)
}
