package guardcmd

import (
	"fmt"
	"guard"

	"github.com/spf13/cobra"
)

var defaultFactory = newRootCmdFactory(func() *cobra.Command {
	return &cobra.Command{
		Use:  "guard",
		Long: `SpiderGuard is a firewall`,
		Example: `$ guard run
$ guard check profile.json`,
		SilenceUsage: true,
		Version:      onlyVersionText(),
	}
})

const fullDocsFooter = `Full documentation is available at: https://spiderguard.yonomesh.com/doc`

func init() {
	defaultFactory.Apply(func(rootCmd *cobra.Command) {
		rootCmd.SetVersionTemplate("{{.Version}}\n")
		rootCmd.SetHelpTemplate(rootCmd.HelpTemplate() + "\n" + fullDocsFooter + "\n")
	})
}

func onlyVersionText() string {
	_, f := guard.Version()
	return f
}

// exitError carries the exit code from CommandFunc to Main()
type exitError struct {
	ExitCode int
	Err      error
}

func (e *exitError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("exiting with code %d", e.ExitCode)
	}
	return e.Err.Error()
}
