package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type Request struct {
	Expression string `json:"expression"`
}

type SuccessResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	errRes := ErrorResponse{Error: err.Error()}
	err = json.NewEncoder(w).Encode(errRes)
	if err != nil {
		slog.Error("Failed to write response", slog.String("error", err.Error()))
	}
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func expectMethodMiddleware(methods ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, m := range methods {
				if r.Method == m {
					next.ServeHTTP(w, r)
					return
				}
			}

			w.WriteHeader(http.StatusMethodNotAllowed)
			writeError(w, fmt.Errorf("expected one of the methods: %s", strings.Join(methods, ", ")))
			return
		})
	}
}

type expressionRequest struct {
	Expression string `json:"expression"`
}

type expressionSimpleResponse struct {
	Id uint64 `json:"id"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	exp := expressionRequest{}
	err := json.NewDecoder(r.Body).Decode(&exp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, invalidRequestBodyError)
		return
	}

	expId, err := orchestrator.CreateExpression(exp.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(expressionSimpleResponse{Id: expId})
}

type expressionsResponse struct {
	Expressions []*expressionResponse `json:"expressions"`
}

func expressionsHandler(w http.ResponseWriter, r *http.Request) {
	expressions, err := orchestrator.GetAllExpressions()
	if err != nil {
		slog.Error("Failed to get expressions", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return
	}

	response := expressionsResponse{
		Expressions: expressions,
	}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		slog.Error("Failed to write expression response", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return
	}
}

func expressionHandler(w http.ResponseWriter, r *http.Request) {
	expressionId, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, invalidIdInUrlError)
		return
	}

	expression, err := orchestrator.GetExpression(expressionId)
	if err != nil {
		if errors.Is(err, expressionNotFoundError) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			slog.Error("Failed to get expression", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}
		writeError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(expression)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
	}
}

type taskRequest struct {
	Id     uint64  `json:"id"`
	Result float64 `json:"result"`
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		task := &taskRequest{}
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeError(w, invalidRequestBodyError)
			return
		}
		err = orchestrator.CompleteTask(task.Id, task.Result)
		if err != nil {
			if errors.Is(err, taskNotFoundError) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				slog.Warn("Client tried to complete a task that is already completed or canceled", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
			}
			writeError(w, err)
			return
		}
		return
	}

	task, err := orchestrator.StartProcessingNextTask()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		writeError(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
	}
}

func (a *Application) RunServer() error {
	http.Handle("/api/v1/calculate", commonMiddleware(
		expectMethodMiddleware("POST")(
			http.HandlerFunc(calculateHandler),
		),
	),
	)
	http.Handle("/api/v1/expressions", commonMiddleware(
		expectMethodMiddleware("GET")(
			http.HandlerFunc(expressionsHandler),
		),
	),
	)
	http.Handle("/api/v1/expressions/{id}", commonMiddleware(
		expectMethodMiddleware("GET")(
			http.HandlerFunc(expressionHandler),
		),
	),
	)

	http.Handle("/internal/task", commonMiddleware(
		expectMethodMiddleware("GET", "POST")(
			http.HandlerFunc(taskHandler),
		),
	),
	)

	slog.Info(fmt.Sprintf("Listening on port %s", a.config.Port))
	return http.ListenAndServe(a.config.Addr+":"+a.config.Port, nil)
}
