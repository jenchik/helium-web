package app

import (
	"context"
	_ "net/http/pprof" // nolint:gosec // Enable profiling

	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/helium/web"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type (
	// Params struct
	Params struct {
		dig.In

		Logger *zap.SugaredLogger
		Server web.Service
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

	s.Logger.Info("Application stopped")

	return nil
}
