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

var TaskEndpoint = "http://localhost:8080/internal/task"
var workerIdCounter = atomic.Uint64{}

func RunWorker() {
	worker := newAgentWorker()
	slog.Info("Started a new agent worker", slog.Uint64("workerId", worker.id))
	for {
		time.Sleep(1 * time.Second)
		go worker.processTask()
	}
}

type agentWorker struct {
	id uint64
}

func newAgentWorker() *agentWorker {
	return &agentWorker{
		id: workerIdCounter.Add(1),
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
		return
	}
	taskResult := &taskRequest{
		Id:     task.Id,
		Result: res,
	}
	err = w.sendTaskResult(taskResult)
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
	request, err := http.NewRequest("GET", TaskEndpoint, nil)
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

type taskRequest struct {
	Id     uint64  `json:"id"`
	Result float64 `json:"result"`
}

func (w *agentWorker) sendTaskResult(task *taskRequest) error {
	client := &http.Client{}

	reqBody := new(bytes.Buffer)
	err := json.NewEncoder(reqBody).Encode(task)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", TaskEndpoint, reqBody)
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
