package app

import (
	"context"
	_ "net/http/pprof" // nolint:gosec // Enable profiling

	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/helium/web"
	"github.com/jenchik/workers"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type (
	// Params struct
	Params struct {
		dig.In

		Logger  *zap.SugaredLogger
		Server  web.Service
		Workers *workers.Group
	}

	serveApp struct {
		*Params
	}
)

func newAppServe(params Params) helium.App {
	return serveApp{
		&params,
	}
}

// Run application
func (s serveApp) Run(ctx context.Context) error {
	s.Logger.Info("Running workers...")
	s.Workers.Run()

	s.Logger.Infow("Run servers")
	if err := s.Server.Start(); err != nil {
		return err
	}

	s.Logger.Info("Startup done")
	<-ctx.Done()

	s.Logger.Info("Stop servers")
	if err := s.Server.Stop(); err != nil {
		return err
	}

	s.Logger.Info("Stopping workers...")
	s.Workers.Stop()

	s.Logger.Info("Waiting workers...")
	err := s.Workers.Wait(context.TODO())

	s.Logger.Info("Application stopped")

	return err
}
