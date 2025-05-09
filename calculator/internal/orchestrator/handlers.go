package orchestrator

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dzherb/go_calculator/internal/auth"
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

func WriteError(w http.ResponseWriter, err error) {
	errRes := ErrorResponse{Error: err.Error()}

	err = json.NewEncoder(w).Encode(errRes)
	if err != nil {
		slog.Error(
			"Failed to write response",
			slog.String("error", err.Error()),
		)
	}
}

type ExpressionRequest struct {
	Expression string `json:"expression"`
}

type ExpressionSimpleResponse struct {
	Id uint64 `json:"id"`
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	exp := ExpressionRequest{}

	err := json.NewDecoder(r.Body).Decode(&exp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, errInvalidRequestBody)

		return
	}

	expId, err := orchestrator.CreateExpression(exp.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		WriteError(w, err)

		return
	}

	slog.Info(
		"Created new expression",
		slog.String("expression", exp.Expression),
		slog.Uint64("expressionId", expId),
	)

	err = json.NewEncoder(w).Encode(ExpressionSimpleResponse{Id: expId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteError(w, err)

		return
	}
}

type ExpressionsResponse struct {
	Expressions []*ExpressionResponse `json:"expressions"`
}

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	expressions, err := orchestrator.GetAllExpressions()
	if err != nil {
		slog.Error(
			"Failed to get expressions",
			slog.String("error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		WriteError(w, err)

		return
	}

	response := ExpressionsResponse{
		Expressions: expressions,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		slog.Error(
			"Failed to write expression response",
			slog.String("error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		WriteError(w, err)

		return
	}
}

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	expressionId, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, errInvalidIdInUrl)

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

		WriteError(w, err)

		return
	}

	err = json.NewEncoder(w).Encode(expression)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteError(w, err)
	}
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var authService = auth.NewService()

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	baseAuthHandler(w, r, authService.Register)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	baseAuthHandler(w, r, authService.Login)
}

func baseAuthHandler(
	w http.ResponseWriter,
	r *http.Request,
	authAction func(username, password string) (auth.AccessPayload, error),
) {
	authReq := authRequest{}

	err := json.NewDecoder(r.Body).Decode(&authReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, errInvalidRequestBody)

		return
	}

	resp, err := authAction(authReq.Username, authReq.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, err)

		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteError(w, err)
	}
}
