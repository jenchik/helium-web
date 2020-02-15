package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jenchik/helium-web/app"
	"github.com/jenchik/helium-web/pkg/core"

	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/helium/grace"
	"github.com/im-kulikov/helium/logger"
	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/helium/settings"
	"github.com/spf13/pflag"
	"go.uber.org/dig"
)

var (
	Name      = "app"
	Version   = "0.1.local"
	BuildTime = ""
	GitHash   = ""
	GitTag    = ""

	configFile *string
)

// nolint:gochecknoinits
func init() {
	if Name == "app" {
		Name = filepath.Base(os.Args[0])
	}
}

func run(mod module.Module) {
	h, err := helium.New(&helium.Settings{
		Name:         Name,
		File:         *configFile,
		BuildTime:    BuildTime,
		BuildVersion: Version,
	}, mod.Append(
		grace.Module,
		settings.Module,
		logger.Module,
		module.New(appInfo),
	))
	helium.Catch(dig.RootCause(err))

	helium.Catch(dig.RootCause(h.Invoke(core.Verbose)))

	helium.Catch(dig.RootCause(h.Run()))
}

func main() {
	version := pflag.Bool("version", false, fmt.Sprintf("show the %s version information", Name))
	configFile = pflag.StringP("config", "c", "config.yml", "path to config file")
	pflag.Parse()

	if *version {
		rtl := fmt.Sprintf("%s (%d/%s/%d)", runtime.Version(), runtime.NumCPU(), runtime.Compiler, runtime.NumCgoCall())
		fmt.Printf(
			"Version: %s (%s)\nBuild: %s\nGit commit hash: %s\nRuntime: %s\n",
			Version,
			GitTag,
			BuildTime,
			GitHash,
			rtl,
		)
		return
	}

	run(app.Module)
}

func appInfo(s *settings.Core) core.AppInfo {
	return core.AppInfo{
		Name:         s.Name,
		BuildTime:    s.BuildTime,
		BuildVersion: s.BuildVersion,
		GitHash:      GitHash,
		GitTag:       GitTag,
	}
}
