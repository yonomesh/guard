package unicmd

import (
	"regexp"
	"sync"

	"github.com/spf13/cobra"
)

// Command represents a subcommand.
// Name, CobraFunc, and Short are required.
type Command struct {
	// The name of the subcommand. Must conform to the
	// format described by the RegisterCommand() godoc.
	// Required.
	Name string

	// Usage is a brief message describing the syntax of
	// the subcommand's flags and args. Use [] to indicate
	// optional parameters and <> to enclose literal values
	// intended to be replaced by the user. Do not prefix
	// the string with "caddy" or the name of the command
	// since these will be prepended for you; only include
	// the actual parameters for this command.
	Usage string

	// Short is a one-line message explaining what the
	// command does. Should not end with punctuation.
	// Required.
	Short string

	// Long is the full help text shown to the user.
	// Will be trimmed of whitespace on both ends before
	// being printed.
	Long string

	// CobraFunc configures the command using Cobra APIs.
	CobraFunc func(*cobra.Command)
}

var (
	commandsMu sync.RWMutex
	commands   = make(map[string]Command)
)

// Commands returns a list of commands initialised by
// RegisterCommand
func Commands() map[string]Command {
	commandsMu.RLock()
	defer commandsMu.RUnlock()

	return commands
}

// CommandFunc is a command's function. It runs the
// command and returns the proper exit code along with
// any error that occurred.
type CommandFunc func(Flags) (int, error)

func init() {
	RegisterCommand(Command{
		Name:  "test",
		Usage: "[--hello <text>]",
		Short: "just test",
		Long:  `Hello Yonomesh`,
		CobraFunc: func(c *cobra.Command) {
			c.Flags().StringP("hello", "", "", "test cmd hello")
			c.RunE = CommandFuncToCobraRunE(cmdTest)
		},
	})

	RegisterCommand(Command{
		Name:  "start",
		Usage: "[--config <path>] [--profile <path> ] [--watch] [--pidfile <file>]",
		Short: "Starts the Guard process in the background and then returns",
		Long: `
Starts the Guard process, optionally bootstrapped with an initial profile and config file.
This command unblocks after the server starts running or fails to run.

On Windows, the spawned child process will remain attached to the terminal, so
closing the window will forcefully stop Guard.
`,
		CobraFunc: func(c *cobra.Command) {
			c.Flags().StringP("config", "c", "", "Configuration file")
			c.Flags().StringP("profile", "p", "", "Profile")
			c.Flags().StringP("pidfile", "", "", "Path of file to which to write process ID")
			c.RunE = CommandFuncToCobraRunE(cmdStart)
		},
	})
}

// RegisterCommand registers the command cmd.
// cmd.Name must be unique and conform to the
// following format:
//
//   - lowercase
//   - alphanumeric and hyphen characters only
//   - cannot start or end with a hyphen
//   - hyphen cannot be adjacent to another hyphen
//
// This function panics if the name is already registered,
// if the name does not meet the described format, or if
// any of the fields are missing from cmd.
//
// This function should be used in init().
func RegisterCommand(cmd Command) {
	commandsMu.Lock()
	defer commandsMu.Unlock()

	if cmd.Name == "" {
		panic("command name is required")
	}
	if cmd.CobraFunc == nil {
		panic("command function missing")
	}
	if cmd.Short == "" {
		panic("command short string is required")
	}
	if _, exists := commands[cmd.Name]; exists {
		panic("command already registered: " + cmd.Name)
	}
	if !commandNameRegex.MatchString(cmd.Name) {
		panic("invalid command name")
	}
	defaultFactory.Apply(func(rootCmd *cobra.Command) {
		rootCmd.AddCommand(UniCmdToCobra(cmd))
	})
	commands[cmd.Name] = cmd
}

var commandNameRegex = regexp.MustCompile(`^[a-z0-9]$|^([a-z0-9]+-?[a-z0-9]*)+[a-z0-9]$`)
