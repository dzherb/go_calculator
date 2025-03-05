package orchestrator

import (
	"fmt"
	"go_calculator/pkg/calculator"
	"log/slog"
	"time"
)

const TaskMaxTimeToLive = 3 * time.Second

type expressionStatus string

const (
	waitingForProcessing expressionStatus = "waiting for processing"
	processing           expressionStatus = "processing"
	processed            expressionStatus = "processed"
	failed               expressionStatus = "failed"
)

type expressionResponse struct {
	Id     uint64           `json:"id"`
	Status expressionStatus `json:"status"`
	Result *float64         `json:"result"`
}

func newExpressionResponse(expression *calculator.Expression) (*expressionResponse, error) {
	var status expressionStatus
	var result *float64

	if expression.IsFailed {
		status = failed
		result = nil
	} else if expression.IsEvaluated() {
		status = processed
		r, err := expression.GetResult()
		if err != nil {
			return nil, err
		}
		result = &r
	} else if expression.IsProcessing {
		status = processing
	} else {
		status = waitingForProcessing
	}

	return &expressionResponse{
		Id:     expression.Id,
		Status: status,
		Result: result,
	}, nil
}

type taskResponse struct {
	Id            uint64  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime uint64  `json:"operation_time"`
}

func newTaskResponse(task *calculator.Task) (*taskResponse, error) {
	arg1, arg2 := task.GetArguments()
	operator := task.GetOperator()
	return &taskResponse{
		Id:            task.Id,
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     operator,
		OperationTime: getOperationTime(operator),
	}, nil
}

func getOperationTime(operator string) uint64 {
	return 1
}

var orchestrator = Orchestrator{
	expressionStorage: ExpressionStorageInstance,
	taskStorage:       TaskStorageInstance,
}

type Orchestrator struct {
	expressionStorage Storage[*calculator.Expression]
	taskStorage       Storage[*calculator.Task]
}

func (o *Orchestrator) CreateExpression(expression string) (uint64, error) {
	exp, err := calculator.NewExpression(expression)
	if err != nil {
		return 0, err
	}
	o.expressionStorage.Put(exp)
	return exp.Id, nil
}

func (o *Orchestrator) GetExpression(id uint64) (*expressionResponse, error) {
	exp, ok := o.expressionStorage.Get(id)
	if !ok {
		return nil, expressionNotFoundError
	}

	return newExpressionResponse(exp)
}

func (o *Orchestrator) GetAllExpressions() ([]*expressionResponse, error) {
	results := make([]*expressionResponse, 0)
	for _, exp := range o.expressionStorage.GetAll() {
		resExp, err := newExpressionResponse(exp)
		if err != nil {
			return nil, err
		}
		results = append(results, resExp)
	}
	return results, nil
}

func (o *Orchestrator) StartProcessingNextTask() (*taskResponse, error) {
	for _, exp := range o.expressionStorage.GetAll() {
		task, ok := exp.GetNextTask()
		if !ok {
			continue
		}
		o.taskStorage.Put(task)

		go func() {
			time.Sleep(TaskMaxTimeToLive)
			err := task.Cancel()
			if err != nil {
				return
			}
			slog.Warn(fmt.Sprintf("Task %d is canceled due to exceeded time to live", task.Id))
		}()

		return newTaskResponse(task)

	}
	return nil, noTasksToProcessError
}

func (o *Orchestrator) CompleteTask(taskId uint64, result float64) error {
	task, ok := o.taskStorage.Get(taskId)
	if !ok {
		return taskNotFoundError
	}
	err := task.Complete(result)
	return err
}

func (o *Orchestrator) CancelTask(id uint64) error {
	task, ok := o.taskStorage.Get(id)
	if !ok {
		return taskNotFoundError
	}
	return task.Cancel()
}

func (o *Orchestrator) OnCalculationFailure(taskId uint64) {
	task, ok := o.taskStorage.Get(taskId)
	if !ok {
		return
	}

	task.GetExpression().MarkAsFailed()
	task.Cancel()
}
