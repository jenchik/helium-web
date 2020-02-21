package chapsuk

import (
	"context"

	prov "github.com/chapsuk/worker"
	"github.com/jenchik/helium-web/pkg/workers"
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
	go a.wg.Stop()
}

func (a *adapter) Wait(ctx context.Context) error {
	stop := make(chan struct{})
	go func() {
		defer close(stop)
		a.wg.Stop()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stop:
	}
	return nil
}

func (adapter) Group(ctx context.Context) workers.Group {
	a := &adapter{
		wg: prov.NewGroup(),
	}
	go func() {
		<-ctx.Done()
		a.wg.Stop()
	}()
	return a
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

	a.wg.Add(w)
	return nil
}
