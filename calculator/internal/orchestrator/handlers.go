package orchestrator

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
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
		slog.Error(
			"Failed to write response",
			slog.String("error", err.Error()),
		)
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
		writeError(w, errInvalidRequestBody)

		return
	}

	expId, err := orchestrator.CreateExpression(exp.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeError(w, err)

		return
	}

	slog.Info(
		"Created new expression",
		slog.String("expression", exp.Expression),
		slog.Uint64("expressionId", expId),
	)

	err = json.NewEncoder(w).Encode(expressionSimpleResponse{Id: expId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)

		return
	}
}

type expressionsResponse struct {
	Expressions []*ExpressionResponse `json:"expressions"`
}

func expressionsHandler(w http.ResponseWriter, r *http.Request) {
	expressions, err := orchestrator.GetAllExpressions()
	if err != nil {
		slog.Error(
			"Failed to get expressions",
			slog.String("error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)

		return
	}

	response := expressionsResponse{
		Expressions: expressions,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		slog.Error(
			"Failed to write expression response",
			slog.String("error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)

		return
	}
}

func expressionHandler(w http.ResponseWriter, r *http.Request) {
	expressionId, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, errInvalidIdInUrl)

		return
	}

	expression, err := orchestrator.GetExpression(expressionId)
	if err != nil {
		if errors.Is(err, errExpressionNotFound) {
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
