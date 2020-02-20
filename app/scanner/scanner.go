package scanner

import (
	"context"

	"github.com/im-kulikov/helium/module"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

var Module = module.Module{
	{Constructor: New},
}

type (
	Config struct {
		TestName string `validate:"required"`
	}

	Options struct {
		dig.In

		Config Config
		Logger *zap.SugaredLogger
	}

	Scanner struct {
		opts *Options
	}
)

// New returns new scanner instance
func New(opts Options) *Scanner {
	opts.Logger = opts.Logger.With("service", "scanner")
	return &Scanner{
		opts: &opts,
	}
}

func (s *Scanner) Job(context.Context) {
	s.opts.Logger.Infow("Run scanner", "my-name", s.opts.Config.TestName)
}
