package app

import (
	"github.com/jenchik/helium-web/app/web/control"
	"github.com/jenchik/helium-web/pkg/web"

	"github.com/im-kulikov/helium/module"
	// "github.com/im-kulikov/helium/web"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

var (
	// Module modules for web-application
	Module = module.New(newAppServe).
		Append(
			module.New(validator.New),
			module.New(configInit),
			web.ServersModule,

			control.Module,
		)
)

type Config struct {
	dig.Out

	Control control.Config
}

func configInit(v *viper.Viper, log *zap.SugaredLogger, validate *validator.Validate) (Config, error) {
	_ = zap.RedirectStdLog(log.With("log_type", "stdout").Desugar())
	log.Infof("Used configuration file: %s", v.ConfigFileUsed())

	var conf Config
	if err := v.Unmarshal(&conf); err != nil {
		return conf, err
	}

	if err := validate.Struct(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}
