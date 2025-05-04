package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"
)

var workerIdCounter = atomic.Uint64{}

const workerStartDelay = 100 * time.Millisecond

func (a *Application) Run() {
	for range a.config.TotalWorkers {
		// Запускаем вокреров с небольшой задержкой,
		// так будем более равномерно обращаться к оркестратору
		go a.runWorker()
		time.Sleep(workerStartDelay)
	}

	waitUntilTermination()
}

const pollingInterval = 400 * time.Millisecond

func (a *Application) runWorker() {
	orchestratorBaseUrl := "http://" +
		a.config.orchestratorHost +
		":" +
		a.config.orchestratorPort

	worker := a.newAgentWorker(orchestratorBaseUrl)

	slog.Info(
		"Started a new agent worker",
		slog.Uint64("workerId", worker.id),
		slog.String("orchestratorUrl", orchestratorBaseUrl),
	)

	for {
		task, err := worker.getTask()
		if err != nil {
			// Не хотим делать подряд неисчислимое количество запросов к оркестратору,
			// если прямо сейчас не смогли получить задачу, поэтому спим
			time.Sleep(pollingInterval)
			continue
		}

		slog.Info(
			"Got a taskToProcess from the orchestrator",
			slog.Uint64("taskId", task.Id),
			slog.Uint64("workerId", worker.id),
		)
		worker.processTask(task)
	}
}

func (a *Application) newAgentWorker(orchestratorBaseUrl string) *agentWorker {
	return &agentWorker{
		id:                  workerIdCounter.Add(1),
		orchestratorTaskUrl: orchestratorBaseUrl + "/internal/taskToProcess",
	}
}

func (w *agentWorker) processTask(task *taskToProcess) {
	res, err := compute(task)
	if err != nil {
		errorReq := &taskErrorRequest{Id: task.Id, Error: err.Error()}

		slog.Info(
			"Failed to compute a taskToProcess",
			slog.Uint64("taskId", task.Id),
			slog.Uint64("workerId", w.id),
		)

		err = w.sendTaskResult(errorReq)
		if err != nil {
			slog.Info(
				"Failed to send taskToProcess result",
				slog.Uint64("taskId", task.Id),
				slog.Uint64("workerId", w.id),
			)

			return
		}

		return
	}

	taskReq := &taskSuccessRequest{
		Id:     task.Id,
		Result: res,
	}

	err = w.sendTaskResult(taskReq)
	if err != nil {
		return
	}

	slog.Info(
		"Successfully processed a taskToProcess",
		slog.Uint64("taskId", task.Id),
		slog.Uint64("workerId", w.id),
	)
}

func (w *agentWorker) getTask() (*taskToProcess, error) {
	client := &http.Client{}

	request, err := http.NewRequest(http.MethodGet, w.orchestratorTaskUrl, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			slog.Error(
				"Failed to close response body",
				slog.String("error", err.Error()),
			)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"failed to retrieve a taskToProcess from the orchestrator (%d status code)",
			response.StatusCode,
		)
	}

	task := &taskToProcess{}

	err = json.NewDecoder(response.Body).Decode(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

type taskProcessed interface{}

type taskSuccessRequest struct {
	Id     uint64  `json:"id"`
	Result float64 `json:"result"`
}

type taskErrorRequest struct {
	Id    uint64 `json:"id"`
	Error string `json:"error"`
}

func (w *agentWorker) sendTaskResult(task taskProcessed) error {
	client := &http.Client{}
	reqBody := new(bytes.Buffer)

	err := json.NewEncoder(reqBody).Encode(task)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		w.orchestratorTaskUrl,
		reqBody,
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"no taskToProcess result confirmation from the orchestrator (%d status code)",
			response.StatusCode,
		)
	}

	return nil
}

func compute(task *taskToProcess) (float64, error) {
	time.Sleep(time.Duration(task.OperationTime))

	var result float64

	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
	case "-":
		result = task.Arg1 - task.Arg2
	case "*":
		result = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}

		result = task.Arg1 / task.Arg2
	default:
		return 0, fmt.Errorf(
			"invalid or unsupported operation '%s'",
			task.Operation,
		)
	}

	return result, nil
}
