package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"
)

var workerIdCounter = atomic.Uint64{}

func (a *Application) Run() {
	for _ = range a.config.TotalWorkers {
		go a.runWorker()
	}
	waitUntilTermination()
}

func (a *Application) runWorker() {
	orchestratorBaseUrl := "http://" + a.config.orchestratorHost + ":" + a.config.orchestratorPort
	worker := a.newAgentWorker(orchestratorBaseUrl)
	slog.Info(
		"Started a new agent worker",
		slog.Uint64("workerId", worker.id),
		slog.String("orchestratorUrl", orchestratorBaseUrl),
	)
	for {
		go worker.processTask()
		time.Sleep(1 * time.Second)
	}
}

type agentWorker struct {
	id                  uint64
	orchestratorTaskUrl string
}

func (a *Application) newAgentWorker(orchestratorBaseUrl string) *agentWorker {
	return &agentWorker{
		id:                  workerIdCounter.Add(1),
		orchestratorTaskUrl: orchestratorBaseUrl + "/internal/task",
	}
}

func (w *agentWorker) processTask() {
	task, err := w.getTask()
	if err != nil {
		return
	}
	slog.Info(
		"Got a task from the orchestrator",
		slog.Uint64("taskId", task.Id),
		slog.Uint64("workerId", w.id),
	)
	res, err := compute(task)
	if err != nil {
		errorReq := &taskErrorRequest{Id: task.Id, Error: err.Error()}
		err = w.sendTaskResult(errorReq)
		slog.Info(
			"Failed to compute a task",
			slog.Uint64("taskId", task.Id),
			slog.Uint64("workerId", w.id),
		)
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
		"Successfully processed a task",
		slog.Uint64("taskId", task.Id),
		slog.Uint64("workerId", w.id),
	)
}

type taskResponse struct {
	Id            uint64  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime uint64  `json:"operation_time"`
}

func (w *agentWorker) getTask() (*taskResponse, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", w.orchestratorTaskUrl, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("couldn't get a task from the orchestrator, received %d status", response.StatusCode))
	}

	task := &taskResponse{}
	err = json.NewDecoder(response.Body).Decode(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

type taskRequest interface{}

type taskSuccessRequest struct {
	Id     uint64  `json:"id"`
	Result float64 `json:"result"`
}

type taskErrorRequest struct {
	Id    uint64 `json:"id"`
	Error string `json:"error"`
}

func (w *agentWorker) sendTaskResult(task taskRequest) error {
	client := &http.Client{}

	reqBody := new(bytes.Buffer)
	err := json.NewEncoder(reqBody).Encode(task)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", w.orchestratorTaskUrl, reqBody)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("got no task result confirmation from the orchestrator, received %d status", response.StatusCode))
	}

	return nil
}

func compute(task *taskResponse) (float64, error) {
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
		return 0, fmt.Errorf("invalid or unsupported operation: %s", task.Operation)
	}

	return result, nil
}
