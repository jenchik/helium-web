package core

import (
	"runtime"

	"go.uber.org/zap"
)

type AppInfo struct {
	Name         string `json:"name"`
	BuildTime    string `json:"build_time"`
	BuildVersion string `json:"build_version"`
	GitHash      string `json:"git_commit_hash"`
	GitTag       string `json:"git_tag"`
}

func Verbose(l *zap.SugaredLogger, a AppInfo) {
	l.Infow("Build",
		"version", a.BuildVersion,
		"build_time", a.BuildTime,
		"git_tag", a.GitTag,
		"git_commit_hash", a.GitHash,
	)
	rtl := []interface{}{
		"cgo_call", runtime.NumCgoCall(),
		"compiler", runtime.Compiler,
		"num_cpu", runtime.NumCPU(),
		"version", runtime.Version(),
	}
	l.Infow("Runtime", rtl...)
}
