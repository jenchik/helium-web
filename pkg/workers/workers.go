package workers

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/dig"
)

type (
	ConfigWorkers map[string]Config

	Config struct {
		Cron        string
		Disabled    bool
		Immediately bool
		Timer       time.Duration
		Ticker      time.Duration
		// Lock        string
	}

	Job func(context.Context)

	// Locker interface
	Locker interface {
		Lock() error
		Unlock()
	}

	// Params is dependencies for create workers slice
	Params struct {
		dig.In

		Context context.Context
		Workers Workers
		Jobs    map[string]Job
		Config  ConfigWorkers `optional:"true"`
		Locker  Locker        `optional:"true"`
	}

	Option struct {
		Name   string
		Job    Job
		Config Config
		Locker Locker
	}

	Group interface {
		Add(Option) error
		Run()
		Stop()
		Wait(context.Context) error
	}

	Workers interface {
		Group(context.Context) Group
	}
)

// NewWorkersGroup returns workers group with injected workers
func NewWorkersGroup(p Params) (Group, error) {
	switch {
	case p.Config == nil || len(p.Config) == 0:
		return nil, ErrEmptyConfig
	case p.Jobs == nil || len(p.Jobs) == 0:
		return nil, ErrEmptyWorkers
	}

	wg := p.Workers.Group(p.Context)
	for name, config := range p.Config {
		job, found := p.Jobs[name]
		if !found || job == nil {
			return nil, fmt.Errorf("%w for '%s'", ErrEmptyJob, name)
		}
		err := wg.Add(Option{
			Config: config,
			Locker: p.Locker,
			Name:   name,
			Job:    job,
		})
		if err != nil {
			// all or nothing
			return nil, fmt.Errorf("%w for '%s'", ErrCreateWorker, name)
		}
	}
	return wg, nil
}
