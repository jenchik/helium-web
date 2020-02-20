package app

import (
	"github.com/jenchik/helium-web/app/scanner"
	"github.com/jenchik/helium-web/pkg/workers"
	"go.uber.org/dig"
)

type JobsParams struct {
	dig.In

	Scanner *scanner.Scanner
}

func jobs(params JobsParams) map[string]workers.Job {
	return map[string]workers.Job{
		"test_scanner": params.Scanner.Job,
	}
}
