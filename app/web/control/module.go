package control

import (
	"net/http"

	"github.com/jenchik/helium-web/pkg/web"

	"github.com/im-kulikov/helium/module"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

// Module of Control API service
var Module = module.Module{
	{
		Constructor: New,
		Options: []dig.ProvideOption{
			dig.Name("control_handler")},
	},
	{Constructor: NewControlServer},
}

type ctrlServerParams struct {
	dig.In

	Handler http.Handler `name:"control_handler"`
	Logger  *zap.SugaredLogger
	Config  Config
}

func NewControlServer(p ctrlServerParams) (web.ServerResult, error) {
	p.Logger = p.Logger.With("web", "control")

	return web.NewHTTPServer(p.Config.HTTP, p.Handler, p.Logger.Desugar())
}
