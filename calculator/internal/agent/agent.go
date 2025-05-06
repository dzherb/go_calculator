package agent

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"sync/atomic"
	"time"
)

var workerIdCounter = atomic.Uint64{}

const workerStartDelay = 100 * time.Millisecond

func New() *Agent {
	return &Agent{
		config: ConfigFromEnv(),
	}
}

func (a *Agent) Run() {
	conn, err := NewOrchestratorConn(a.config)
	if err != nil {
		slog.Error(
			"failed to connect to the orchestrator",
			"error", err,
		)

		return
	}

	defer func(conn *OrchestratorConn) {
		err = conn.Close()
		if err != nil {
			slog.Error(
				"failed to close the connection",
				"error", err,
			)
		}
	}(conn)

	a.client = OrchestratorClient(conn)

	a.run()
}

func (a *Agent) run() {
	for range a.config.TotalWorkers {
		// Запускаем вокреров с небольшой задержкой,
		// так будем более равномерно обращаться к оркестратору
		go a.runWorker()
		time.Sleep(workerStartDelay)
	}

	waitUntilTermination()
}

const pollingInterval = 400 * time.Millisecond

func (a *Agent) runWorker() {
	worker := a.newAgentWorker()

	slog.Info(
		"started a new agent worker",
		slog.Uint64("workerId", worker.id),
	)

	for {
		task, err := worker.getTask()
		if err != nil {
			st, ok := status.FromError(err)
			if !ok || st.Code() != codes.ResourceExhausted {
				slog.Error(
					"unexpected error while getting a task",
					"error", err,
					"taskId", task.Id,
					"workerId", worker.id,
				)
			}

			// No tasks available, sleep for a while
			time.Sleep(pollingInterval)
			continue
		}

		slog.Info(
			"got a taskToProcess from the orchestrator",
			"taskId", task.Id,
			"workerId", worker.id,
		)
		worker.processTask(task)
	}
}

func (a *Agent) newAgentWorker() *agentWorker {
	return &agentWorker{
		id:    workerIdCounter.Add(1),
		agent: a,
	}
}
