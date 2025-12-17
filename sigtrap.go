package guard

// Exit codes. Generally, you should NOT
// automatically restart the process if the
// exit code is ExitCodeFailedStartup (1).
const (
	ExitCodeSuccess = iota
	ExitCodeFailedStartup
	ExitCodeForceQuit
	ExitCodeFailedQuit
)
