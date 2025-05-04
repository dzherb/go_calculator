package orchestrator

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dzherb/go_calculator/pkg/calculator"
)

type Orchestrator struct {
	app               *Application
	expressionStorage Storage[*calculator.Expression]
	taskStorage       Storage[*calculator.Task]
}

var orchestrator = Orchestrator{
	expressionStorage: ExpressionStorageInstance,
	taskStorage:       TaskStorageInstance,
}

func (o *Orchestrator) CreateExpression(expression string) (uint64, error) {
	exp, err := calculator.NewExpression(expression)
	if err != nil {
		return 0, err
	}

	o.expressionStorage.Put(exp)

	return exp.Id, nil
}

func (o *Orchestrator) GetExpression(id uint64) (*ExpressionResponse, error) {
	exp, ok := o.expressionStorage.Get(id)
	if !ok {
		return nil, expressionNotFoundError
	}

	return newExpressionResponse(exp)
}

func (o *Orchestrator) GetAllExpressions() ([]*ExpressionResponse, error) {
	results := make([]*ExpressionResponse, 0)

	for _, exp := range o.expressionStorage.GetAll() {
		resExp, err := newExpressionResponse(exp)
		if err != nil {
			return nil, err
		}

		results = append(results, resExp)
	}

	return results, nil
}

func (o *Orchestrator) StartProcessingNextTask() (*TaskResponse, error) {
	for _, exp := range o.expressionStorage.GetAll() {
		task, ok := exp.GetNextTask()
		if !ok {
			continue
		}

		o.taskStorage.Put(task)

		go func() {
			time.Sleep(o.app.config.TaskMaxProcessTime)

			err := task.Cancel()
			if err != nil {
				return
			}

			slog.Warn(
				fmt.Sprintf(
					"Task %d is canceled due to exceeded time to live",
					task.Id,
				),
			)
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

type expressionStatus string

const (
	waitingForProcessing expressionStatus = "waiting for processing"
	processing           expressionStatus = "processing"
	processed            expressionStatus = "processed"
	failed               expressionStatus = "failed"
)

type ExpressionResponse struct {
	Id     uint64           `json:"id"`
	Status expressionStatus `json:"status"`
	Result *float64         `json:"result"`
}

func newExpressionResponse(
	expression *calculator.Expression,
) (*ExpressionResponse, error) {
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

	return &ExpressionResponse{
		Id:     expression.Id,
		Status: status,
		Result: result,
	}, nil
}

type TaskResponse struct {
	Id            uint64        `json:"id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
}

func newTaskResponse(task *calculator.Task) (*TaskResponse, error) {
	arg1, arg2 := task.GetArguments()
	operator := task.GetOperator()

	return &TaskResponse{
		Id:            task.Id,
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     operator,
		OperationTime: orchestrator.getOperationTime(operator),
	}, nil
}

func (o *Orchestrator) getOperationTime(operator string) time.Duration {
	switch operator {
	case "+":
		return o.app.config.AdditionTime
	case "-":
		return o.app.config.DivisionTime
	case "*":
		return o.app.config.MultiplicationTime
	case "/":
		return o.app.config.DivisionTime
	}

	return 0
}
