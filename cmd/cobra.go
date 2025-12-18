package guardcmd

import (
	"fmt"
	"guard"

	"github.com/yonomesh/cobra"
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

func GuardCmdToCobra(guardCmd Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   guardCmd.Name + " " + guardCmd.Usage,
		Short: guardCmd.Short,
		Long:  guardCmd.Long,
	}

	guardCmd.CobraFunc(cmd)

	return cmd
}

// CommandFuncToCobraRunE wraps a Guard CommandFunc for use
// in a cobra command's RunE field.
func CommandFuncToCobraRunE(f CommandFunc) func(cmd *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		status, err := f(Flags{cmd.Flags()}) // key point
		if status > 1 {
			cmd.SilenceErrors = true
			return &exitError{ExitCode: status, Err: err}
		}
		return err
	}
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

func onlyVersionText() string {
	_, f := guard.Version()
	return f
}
