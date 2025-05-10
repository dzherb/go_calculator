package orchestrator

import (
	"context"
	"log/slog"
	"time"

	"github.com/dzherb/go_calculator/internal/repository"
)

type ExpressionsDaemon struct {
}

func NewExpressionsDaemon() *ExpressionsDaemon {
	return &ExpressionsDaemon{}
}

const ExpressionAbortPeriod = time.Minute * 5

func (ed *ExpressionsDaemon) Start(ctx context.Context) {
	slog.Info("starting expressions daemon")

	ticker := time.NewTicker(ExpressionAbortPeriod)

	for {
		select {
		case <-ctx.Done():
			slog.Info("stopping expressions daemon")
			return
		case <-ticker.C:
			go ed.AbortUnprocessedExpressions()
		}
	}
}

func (ed *ExpressionsDaemon) AbortUnprocessedExpressions() {
	expressions, err := ExpressionRepo().Unprocessed()
	if err != nil {
		slog.Error("Failed to retrieve old expressions", "error", err)
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
	}
}
