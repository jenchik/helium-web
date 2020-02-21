package wrks

import (
	"context"

	"github.com/jenchik/helium-web/pkg/workers"
	prov "github.com/jenchik/workers"
)

type (
	adapter struct {
		wg *prov.Group
	}
)

func New() workers.Workers {
	return adapter{}
}

func (a *adapter) Run() {
	a.wg.Run()
}

func (a *adapter) Stop() {
	a.wg.Stop()
}

func (a *adapter) Wait(ctx context.Context) error {
	return a.wg.Wait(ctx)
}

func (adapter) Group(ctx context.Context) workers.Group {
	return &adapter{
		wg: prov.NewGroup(ctx),
	}
}

func jobFunc(f workers.Job) func(context.Context) {
	return f
}

func (a *adapter) Add(opt workers.Option) error {
	if opt.Config.Disabled {
		return nil
	}

	w := prov.New(jobFunc(opt.Job))

	if opt.Config.Timer > 0 {
		w = w.ByTimer(opt.Config.Timer)
	}
	if opt.Config.Ticker > 0 {
		w = w.ByTicker(opt.Config.Ticker)
	}
	if len(opt.Config.Cron) > 0 {
		w = w.ByCronSpec(opt.Config.Cron)
	}
	if opt.Config.Immediately {
		w = w.SetImmediately(true)
	}

	// TODO
	/*
		if opts.Viper.IsSet(key + ".lock") {
			if opts.Locker == nil {
				return nil, errors.Wrap(ErrEmptyLocker, key)
			} else if l, ok := opts.Locker.(LockerSettings); ok {
				locker, err := l.Apply(key, opts.Viper)
				if err != nil {
					return nil, err
				}

				w = w.WithLock(locker)
			} else {
				w = w.WithLock(opts.Locker)
			}
		}
	*/

	return a.wg.Add(w)
}
