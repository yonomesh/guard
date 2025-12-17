package guardcmd

import "runtime/debug"

type moduleInfo struct {
	guardModuleID string
	goModule      *debug.Module
	err           error
}
