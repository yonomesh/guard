package unicmd

import (
	"errors"
	"fmt"
	"os"

	"uni"
)

// Main implements the main function of the caddy command.
// Call this if Guard is to be the main() of your program.
func Main() {
	if len(os.Args) == 0 {
		fmt.Println("[fatal] no arguments provided by OS; args[0] must be command")
		os.Exit(uni.ExitCodeFailedStartup)
	}

	if err := defaultFactory.Build().Execute(); err != nil {
		var exitError *exitError
		if errors.As(err, &exitError) {
			os.Exit(exitError.ExitCode)
		}
		os.Exit(1)
	}
}
