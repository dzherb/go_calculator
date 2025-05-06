package agent

import (
	pb "github.com/dzherb/go_calculator/internal/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrchestratorConn struct {
	conn    *grpc.ClientConn
	closeFn func() error
}

func (c *OrchestratorConn) Close() error {
	return c.closeFn()
}

func NewOrchestratorConn(cfg *Config) (*OrchestratorConn, error) {
	addr := cfg.orchestratorHost + ":" + cfg.orchestratorPort
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	conn.Connect()

	return &OrchestratorConn{conn: conn, closeFn: conn.Close}, nil
}

func OrchestratorClient(conn *OrchestratorConn) pb.TaskServiceClient {
	return pb.NewTaskServiceClient(conn.conn)
}
