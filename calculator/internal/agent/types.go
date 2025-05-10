package agent

import pb "github.com/dzherb/go_calculator/calculator/internal/gen"

type Agent struct {
	config *Config
	client pb.TaskServiceClient
}

type agentWorker struct {
	id    uint64
	agent *Agent
}
