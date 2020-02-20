package workers

import (
	"fmt"
)

type errBase string

func (e errBase) Error() string {
	return string(e)
}

var (
	ErrWrongConfig = errBase("configure error")

	// ErrMissingKey when config key for worker is missing
	ErrMissingKey = fmt.Errorf("%w; missing worker key", ErrWrongConfig)

	// ErrEmptyConfig when viper not passed to params
	ErrEmptyConfig = fmt.Errorf("%w; empty config", ErrWrongConfig)

	// ErrEmptyWorkers when workers not passed to params
	ErrEmptyWorkers = fmt.Errorf("%w; empty workers", ErrWrongConfig)

	// ErrEmptyLocker when locker required,
	// but not passed to params
	ErrEmptyLocker = fmt.Errorf("%w; empty locker", ErrWrongConfig)

	// ErrEmptyJob when worker job is nil
	ErrEmptyJob = fmt.Errorf("%w; empty job", ErrWrongConfig)
)
