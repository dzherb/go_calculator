package agent

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/dzherb/go_calculator/internal/gen"
)

func (w *agentWorker) client() pb.TaskServiceClient {
	return w.agent.client
}

func (w *agentWorker) processTask(task *pb.TaskToProcess) {
	res, err := compute(task)
	resResp := &pb.TaskResult{
		Id:     task.Id,
		Result: res,
	}

	if err != nil {
		resResp.Error = err.Error()

		slog.Info(
			"Failed to compute a taskToProcess",
			"taskId", task.Id,
			"workerId", w.id,
		)

		err = w.sendTaskResult(resResp)
		if err != nil {
			slog.Info(
				"failed to send task result",
				"taskId", task.Id,
				"workerId", w.id,
			)

			return
		}

		return
	}

	err = w.sendTaskResult(resResp)
	if err != nil {
		slog.Info(
			"failed to send task result",
			"taskId", task.Id,
			"workerId", w.id,
		)

		return
	}

	slog.Info(
		"successfully processed a task",
		"taskId", task.Id,
		"workerId", w.id,
		"res", res,
	)
}

func (w *agentWorker) getTask() (*pb.TaskToProcess, error) {
	client := w.client()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	task, err := client.GetTask(ctx, &pb.GetTaskRequest{})
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (w *agentWorker) sendTaskResult(task *pb.TaskResult) error {
	client := w.client()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := client.AddResult(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func compute(task *pb.TaskToProcess) (float64, error) {
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
