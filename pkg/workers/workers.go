package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/jenchik/workers"
	// "github.com/spf13/viper"
	"go.uber.org/dig"
)

type (
	Job = workers.Job

	ConfigWorkers map[string]Config

	Config struct {
		Cron        string
		Disabled    bool
		Immediately bool
		Timer       time.Duration
		Ticker      time.Duration
		// Lock        string
	}

	// Params is dependencies for create workers slice
	Params struct {
		dig.In

		Jobs   map[string]Job
		Config ConfigWorkers  `optional:"true"`
		Locker workers.Locker `optional:"true"`
	}

	// LockerSettings creates copy of locker and applies settings
	// LockerSettings interface {
	// 	Apply(key string, v *viper.Viper) (workers.Locker, error)
	// }

	options struct {
		Name   string
		Job    Job
		Config Config
		Locker workers.Locker
	}
)

func nopJob(_ context.Context) {}

// NewWorkersGroup returns workers group with injected workers
func NewWorkersGroup(ctx context.Context, wrks []*workers.Worker) (*workers.Group, error) {
	var items = make([]*workers.Worker, 0, len(wrks))

	for i := range wrks {
		if wrks[i] != nil {
			items = append(items, wrks[i])
		}
	}

	wg := workers.NewGroup(ctx)
	return wg, wg.Add(items...)
}

// NewWorkers returns wrapped workers slice created by config settings
func NewWorkers(p Params) ([]*workers.Worker, error) {
	switch {
	case p.Config == nil || len(p.Config) == 0:
		return nil, ErrEmptyConfig
	case p.Jobs == nil || len(p.Jobs) == 0:
		return nil, ErrEmptyWorkers
	}

	workers := make([]*workers.Worker, 0, len(p.Jobs))
	for name, config := range p.Config {
		job, found := p.Jobs[name]
		if !found || job == nil {
			return nil, fmt.Errorf("%w for '%s'", ErrEmptyJob, name)
		}
		wrk, err := workerByConfig(options{
			Config: config,
			Locker: p.Locker,
			Name:   name,
			Job:    job,
		})
		if err != nil {
			// all or nothing
			return nil, err
		}
		workers = append(workers, wrk)
	}
	return workers, nil
}

func workerByConfig(opts options) (*workers.Worker, error) {
	if opts.Config.Disabled {
		return workers.New(nopJob), nil
	}

	w := workers.New(opts.Job)

	if opts.Config.Timer > 0 {
		w = w.ByTimer(opts.Config.Timer)
	}
	if opts.Config.Ticker > 0 {
		w = w.ByTicker(opts.Config.Ticker)
	}
	if len(opts.Config.Cron) > 0 {
		w = w.ByCronSpec(opts.Config.Cron)
	}
	if opts.Config.Immediately {
		w = w.SetImmediately(true)
	}

	// if opts.Viper.IsSet(key + ".lock") {
	// 	if opts.Locker == nil {
	// 		return nil, errors.Wrap(ErrEmptyLocker, key)
	// 	} else if l, ok := opts.Locker.(LockerSettings); ok {
	// 		locker, err := l.Apply(key, opts.Viper)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		w = w.WithLock(locker)
	// 	} else {
	// 		w = w.WithLock(opts.Locker)
	// 	}
	// }

	return w, nil
}
