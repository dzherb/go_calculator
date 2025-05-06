package orchestrator

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net"
	"strconv"

	pb "github.com/dzherb/go_calculator/internal/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Application) ServeGRPC() {
	addr := a.config.Host + ":" + a.config.GRPCPort

	lis, err := net.Listen("tcp", a.config.Host+":"+a.config.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterTaskServiceServer(s, &grpcServer{})

	slog.Info("GRPC server is listening on " + addr)

	if err = s.Serve(lis); err != nil {
		slog.Error("failed to start grpc server", "error", err)
	}
}

type grpcServer struct {
	pb.UnimplementedTaskServiceServer
}

func (gs *grpcServer) GetTask(
	_ context.Context,
	_ *pb.GetTaskRequest,
) (*pb.TaskToProcess, error) {
	task, err := orchestrator.StartProcessingNextTask()
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	return task, nil
}

func (gs *grpcServer) AddResult(
	_ context.Context,
	task *pb.TaskResult,
) (*pb.AddResultResponse, error) {
	if task.Error != "" {
		slog.Error(
			"Agent returned calculation error",
			slog.String("error", task.Error),
		)

		err := orchestrator.OnCalculationFailure(task.Id)
		if err != nil {
			slog.Error(err.Error())
		}

		return nil, nil
	}

	slog.Info(
		"Got task result",
		slog.String("id", strconv.FormatUint(task.Id, 10)),
	)

	err := orchestrator.CompleteTask(task.Id, task.Result)
	if err != nil {
		if !errors.Is(err, errTaskNotFound) {
			slog.Warn(
				"Agent tried to complete a task that is already completed or canceled",
				"error",
				err,
			)
		}

		return nil, status.Error(codes.NotFound, err.Error())
	}

	return nil, nil
}
