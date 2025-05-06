package orchestrator

import "context"

type App struct {
	config *Config
}

func New() *App {
	app := &App{
		config: ConfigFromEnv(),
	}
	orchestrator.app = app

	return app
}

func (a *App) Serve() error {
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func() {
		errChan <- a.ServeGRPC(ctx)
	}()

	go func() {
		errChan <- a.ServeHTTP(ctx)
	}()

	return <-errChan
}
