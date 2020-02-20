package workers

import (
	"github.com/im-kulikov/helium/module"
)

// Module of workers
var Module = module.Module{
	{Constructor: NewWorkersGroup},
	{Constructor: NewWorkers},
}
