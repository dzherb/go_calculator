package application

import (
	"encoding/json"
	"fmt"
	"go_calculator/pkg/calculator"
	"log/slog"
	"net/http"
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

func CalculatorHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeError(w, fmt.Errorf("expected POST method"))
		return
	}

	request := new(Request)

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	if request.Expression == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, fmt.Errorf("expression is required"))
		return
	}

	result, err := calculator.Calculate(request.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeError(w, fmt.Errorf("expression is not valid"))
		return
	}

	err = json.NewEncoder(w).Encode(SuccessResponse{Result: result})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, fmt.Errorf("internal server error"))
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalculatorHandler)

	slog.Info(fmt.Sprintf("Listening on port %s", a.config.Port))
	return http.ListenAndServe(a.config.Addr+":"+a.config.Port, nil)
}
