package orchestrator

import (
	"fmt"
	"log/slog"
	"time"

	pb "github.com/dzherb/go_calculator/internal/gen"
	"github.com/dzherb/go_calculator/internal/repository"
	"github.com/dzherb/go_calculator/pkg/calculator"
)

type Orchestrator struct {
	app            *App
	exprMemStorage Storage[*calc.Expression]
	taskMemStorage Storage[*calc.Task]
}

var orchestrator = Orchestrator{
	exprMemStorage: ExpressionStorageInstance,
	taskMemStorage: TaskStorageInstance,
}

var expressionRepo repo.ExpressionRepository

func ExpressionRepo() repo.ExpressionRepository {
	if expressionRepo == nil {
		expressionRepo = repo.NewExpressionRepository()
	}

	return expressionRepo
}

func (o *Orchestrator) CreateExpression(
	expression string,
	userID uint64,
) (uint64, error) {
	expr, err := calc.NewExpression(expression)
	if err != nil {
		return 0, err
	}

	exprFromDB, err := ExpressionRepo().Create(repo.Expression{
		UserID:     userID,
		Expression: expression,
	})
	if err != nil {
		return 0, err
	}

	expr.Id = exprFromDB.ID

	o.exprMemStorage.Put(expr)

	return expr.Id, nil
}

func (o *Orchestrator) GetExpression(id uint64) (repo.Expression, error) {
	return ExpressionRepo().Get(id)
}

func (o *Orchestrator) GetUserExpressions(
	userID uint64,
) ([]repo.Expression, error) {
	return ExpressionRepo().GetForUser(userID)
}

func (o *Orchestrator) StartProcessingNextTask() (*pb.TaskToProcess, error) {
	for expr := range o.exprMemStorage.GetAll() {
		task, ok := expr.GetNextTask()
		if !ok {
			continue
		}

		_, err := ExpressionRepo().Update(repo.Expression{
			ID:     expr.Id,
			Status: repo.ExpressionProcessing,
		})
		if err != nil {
			return nil, err
		}

		o.taskMemStorage.Put(task)

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

		return newTaskToProcess(task)
	}

	return nil, errNoTasksToProcess
}

func (o *Orchestrator) CompleteTask(taskId uint64, result float64) error {
	task, ok := o.taskMemStorage.Get(taskId)
	if !ok {
		return errTaskNotFound
	}

	err := task.Complete(result)
	if err != nil {
		return err
	}

	expr := task.GetExpression()

	if !expr.IsEvaluated() {
		return nil
	}

	res, err := expr.GetResult()
	if err != nil {
		return err
	}

	_, err = ExpressionRepo().Update(repo.Expression{
		ID:     task.GetExpression().Id,
		Status: repo.ExpressionSucceed,
		Result: &res,
	})

	return err
}

func (o *Orchestrator) CancelTask(id uint64) error {
	task, ok := o.taskMemStorage.Get(id)
	if !ok {
		return errTaskNotFound
	}

	return task.Cancel()
}

func (o *Orchestrator) OnCalculationFailure(taskId uint64) error {
	task, ok := o.taskMemStorage.Get(taskId)
	if !ok {
		return errTaskNotFound
	}

	task.GetExpression().MarkAsFailed()

	_, err := ExpressionRepo().Update(repo.Expression{
		ID:     task.GetExpression().Id,
		Status: repo.ExpressionFailed,
	})
	if err != nil {
		return err
	}

	return task.Cancel()
}

type ExpressionResponse struct {
	Id     uint64                `json:"id"`
	Status repo.ExpressionStatus `json:"status"`
	Result *float64              `json:"result"`
}

func newTaskToProcess(
	task *calc.Task,
) (*pb.TaskToProcess, error) { //nolint:unparam
	arg1, arg2 := task.GetArguments()
	operator := task.GetOperator()

	return &pb.TaskToProcess{
		Id:        task.Id,
		Arg1:      arg1,
		Arg2:      arg2,
		Operation: operator,
		OperationTime: uint32( //nolint:gosec
			orchestrator.getOperationTime(operator),
		),
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
