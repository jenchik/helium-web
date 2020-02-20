package control

import (
	"net/http"

	"github.com/jenchik/helium-web/pkg/web"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

type (
	Config struct {
		HTTP web.Config

		Debug    bool
		TestName string
	}

	Options struct {
		dig.In

		Config Config
		Logger *zap.SugaredLogger
	}
)

func New(opt Options) (http.Handler, error) {
	opt.Logger = opt.Logger.With("web", "control")

	mux := http.NewServeMux()
	mux.Handle("/check", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello " + opt.Config.TestName))
	}))

	return mux, nil
}
