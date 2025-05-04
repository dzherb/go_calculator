package agent

type taskToProcess struct {
	Id            uint64  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime uint32  `json:"operation_time"`
}

type agentWorker struct {
	id                  uint64
	orchestratorTaskUrl string
}
