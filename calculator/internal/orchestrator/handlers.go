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
		slog.Error("Failed to write response", slog.String("error", err.Error()))
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
	slog.Info(
		"Created new expression",
		slog.String("expression", exp.Expression),
		slog.Uint64("expressionId", expId),
	)
}

type expressionsResponse struct {
	Expressions []*ExpressionResponse `json:"expressions"`
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
	Error  *string `json:"error"`
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

		if task.Error != nil {
			slog.Error("Agent returned calculation error", slog.String("error", *task.Error))
			orchestrator.OnCalculationFailure(task.Id)
			return
		}

		err = orchestrator.CompleteTask(task.Id, task.Result)
		slog.Info("Got task result", slog.String("id", strconv.FormatUint(task.Id, 10)))
		if err != nil {
			if errors.Is(err, taskNotFoundError) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				slog.Warn("Agent tried to complete a task that is already completed or canceled", slog.String("error", err.Error()))
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
	slog.Info("Orchestrator delegated a task", slog.String("id", strconv.FormatUint(task.Id, 10)))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
	}
}
