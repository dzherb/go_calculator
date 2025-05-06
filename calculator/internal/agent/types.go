package agent

import pb "github.com/dzherb/go_calculator/internal/gen"

type Agent struct {
	config *Config
	client pb.TaskServiceClient
}

type taskToProcess struct {
	Id            uint64  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime uint32  `json:"operation_time"`
}

type agentWorker struct {
	id    uint64
	agent *Agent
}

type taskProcessed interface{}

type taskSuccessRequest struct {
	Id     uint64  `json:"id"`
	Result float64 `json:"result"`
}

type taskErrorRequest struct {
	Id    uint64 `json:"id"`
	Error string `json:"error"`
}
