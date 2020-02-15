package web

import (
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/helium/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type (
	ServerResult = web.ServerResult

	Config struct {
		Address  string
		Disabled bool

		MaxHeaderBytes int
		Network        string
		SkipErrors     bool

		ShutdownTimeout   time.Duration
		ReadTimeout       time.Duration
		ReadHeaderTimeout time.Duration
		WriteTimeout      time.Duration
		IdleTimeout       time.Duration
	}

	// APIParams struct
	APIParams struct {
		dig.In

		Config  Config       `optional:"true"`
		Logger  *zap.Logger  `optional:"true"`
		Handler http.Handler `optional:"true"`
	}

	// MultiServerParams struct
	MultiServerParams struct {
		dig.In

		Logger  *zap.Logger   `optional:"true"`
		Servers []web.Service `group:"services"`
	}

	profileParams struct {
		dig.In

		Logger  *zap.Logger  `optional:"true"`
		Config  Config       `optional:"true"`
		Handler http.Handler `name:"profile_handler" optional:"true"`
	}

	metricParams struct {
		dig.In

		Logger  *zap.Logger  `optional:"true"`
		Config  Config       `optional:"true"`
		Handler http.Handler `name:"metric_handler" optional:"true"`
	}
)

var (
	// ServersModule of web base structs
	ServersModule = module.Module{
		{Constructor: newProfileServer},
		{Constructor: newMetricServer},
		{Constructor: NewAPIServer},
		{Constructor: NewMultiServer},
	}
)

// NewMultiServer returns new multi servers group
func NewMultiServer(p MultiServerParams) (web.Service, error) {
	if p.Logger == nil {
		p.Logger = zap.NewNop()
	}
	return web.New(p.Logger, p.Servers...)
}

func newProfileServer(p profileParams) (ServerResult, error) {
	if p.Logger == nil {
		p.Logger = zap.NewNop()
	}
	p.Logger = p.Logger.With(zap.String("web", "pprof"))
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	if p.Handler != nil {
		mux.Handle("/", p.Handler)
	}
	return NewHTTPServer(p.Config, mux, p.Logger)
}

func newMetricServer(p metricParams) (ServerResult, error) {
	if p.Logger == nil {
		p.Logger = zap.NewNop()
	}
	p.Logger = p.Logger.With(zap.String("web", "metrics"))
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	if p.Handler != nil {
		mux.Handle("/", p.Handler)
	}
	return NewHTTPServer(p.Config, mux, p.Logger)
}

// NewAPIServer creates api server by http.Handler from DI container
func NewAPIServer(p APIParams) (ServerResult, error) {
	if p.Logger == nil {
		p.Logger = zap.NewNop()
	}
	p.Logger = p.Logger.With(zap.String("web", "api"))
	return NewHTTPServer(p.Config, p.Handler, p.Logger)
}

// NewHTTPServer creates http-server that will be embedded into multi-server
func NewHTTPServer(conf Config, h http.Handler, l *zap.Logger) (web.ServerResult, error) {
	var result web.ServerResult

	switch {
	case l == nil:
		return result, web.ErrEmptyLogger
	case h == nil:
		l.Info("Empty handler, skip")
		return result, nil
	case conf.Disabled:
		l.Info("Server disabled")
		return result, nil
	case len(conf.Address) == 0:
		l.Info("Empty bind address, skip")
		return result, nil
	}

	options := []web.HTTPOption{
		web.HTTPListenAddress(conf.Address),
		web.HTTPShutdownTimeout(conf.ShutdownTimeout),
	}

	if len(conf.Network) > 0 {
		options = append(options, web.HTTPListenNetwork(conf.Network))
	}

	if conf.SkipErrors {
		options = append(options, web.HTTPSkipErrors())
	}

	serve, err := web.NewHTTPService(
		&http.Server{
			Handler:           h,
			Addr:              conf.Address,
			ReadTimeout:       conf.ReadTimeout,
			ReadHeaderTimeout: conf.ReadHeaderTimeout,
			WriteTimeout:      conf.WriteTimeout,
			IdleTimeout:       conf.IdleTimeout,
			MaxHeaderBytes:    conf.MaxHeaderBytes,
		},
		options...,
	)

	l.Info("Creates http server",
		zap.String("address", conf.Address))

	result.Server = serve
	return result, err
}
