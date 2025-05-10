package orchestrator

import (
	"context"

	"github.com/dzherb/go_calculator/pkg/security"
)

type App struct {
	config *Config
}

func New() *App {
	cfg := ConfigFromEnv()

	security.Init(security.Config{
		SecretKey:      cfg.SecretKey,
		AccessTokenTTL: cfg.AccessTokenTTL,
	})

	app := &App{
		config: cfg,
	}
	orchestrator.app = app

	return app
}

func (a *App) Serve() error {
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go NewExpressionsDaemon().Start(ctx)

	go func() {
		errChan <- a.ServeGRPC(ctx)
	}()

	go func() {
		errChan <- a.ServeHTTP(ctx)
	}()

	return <-errChan
}
