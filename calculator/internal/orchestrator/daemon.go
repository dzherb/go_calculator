package orchestrator

import (
	"context"
	"log/slog"
	"time"

	"github.com/dzherb/go_calculator/internal/repository"
)

type Daemon struct {
}

func NewDaemon() *Daemon {
	return &Daemon{}
}

const ExprAbortPeriod = time.Minute * 5
const ExprCleanPeriod = time.Minute * 2
const TasksCleanPeriod = time.Minute * 1

func (d *Daemon) Start(ctx context.Context) {
	slog.Info("starting orchestrator daemon")

	ExprAbortTicker := time.Tick(ExprAbortPeriod)
	ExprCleanTicker := time.Tick(ExprCleanPeriod)
	TasksCleanTicker := time.Tick(TasksCleanPeriod)

	for {
		select {
		case <-ctx.Done():
			slog.Info("stopping orchestrator daemon")
			return
		case <-ExprAbortTicker:
			go d.AbortUnprocessedExpr()
		case <-ExprCleanTicker:
			go d.CleanExprStorage()
		case <-TasksCleanTicker:
			go d.CleanTasksStorage()
		}
	}
}

func (d *Daemon) AbortUnprocessedExpr() {
	expressions, err := ExpressionRepo().Unprocessed()
	if err != nil {
		slog.Error("Failed to retrieve unprocessed expressions", "error", err)
		return
	}

	for _, expr := range expressions {
		_, err = ExpressionRepo().Update(repo.Expression{
			ID:     expr.ID,
			Status: repo.ExpressionAborted,
		})
		if err != nil {
			slog.Error("Failed to update expression status",
				"expression_id", expr.ID,
				"error", err,
			)

			return
		}

		orchestrator.exprMemStorage.Delete(expr.ID)
	}
}

func (d *Daemon) CleanExprStorage() {
	var exprs []uint64

	for e := range orchestrator.exprMemStorage.All() {
		if e.IsFailed || e.IsEvaluated() {
			exprs = append(exprs, e.Id)
		}
	}

	for _, id := range exprs {
		orchestrator.exprMemStorage.Delete(id)
	}
}

func (d *Daemon) CleanTasksStorage() {
	var tasks []uint64

	for t := range orchestrator.taskMemStorage.All() {
		if t.IsCanceled || t.IsCompleted {
			tasks = append(tasks, t.Id)
		}
	}

	for _, id := range tasks {
		orchestrator.taskMemStorage.Delete(id)
	}
}
