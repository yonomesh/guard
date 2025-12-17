package guardcmd

import (
	"flag"
	"regexp"
	"sync"

	"github.com/spf13/cobra"
)

// Command represents a subcommand.
// Name, Func, and Short are required.
type Command struct {
	// The name of the subcommand. Must conform to the
	// format described by the RegisterCommand() godoc.
	// Required.
	Name string

	// Usage is a brief message describing the syntax of
	// the subcommand's flags and args. Use [] to indicate
	// optional parameters and <> to enclose literal values
	// intended to be replaced by the user. Do not prefix
	// the string with "guard" or the name of the command
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

	// Flags is the flagset for command.
	// This is ignored if CobraFunc is set.
	Flags *flag.FlagSet

	// Func is a function that executes a subcommand using
	// the parsed flags. It returns an exit code and any
	// associated error.
	// Required if CobraFunc is not set.
	Func CommandFunc

	// CobraFunc allows further configuration of the command
	// via cobra's APIs. If this is set, then Func and Flags
	// are ignored, with the assumption that they are set in
	// this function. A caddycmd.WrapCommandFuncForCobra helper
	// exists to simplify porting CommandFunc to Cobra's RunE.
	CobraFunc func(*cobra.Command)
}

// Commands returns a list of commands initialised by
// RegisterCommand
func Commands() map[string]Command {
	commandsMu.RLock()
	defer commandsMu.RUnlock()

	return commands
}

var (
	commandsMu sync.RWMutex
	commands   = make(map[string]Command)
)

// CommandFunc is a command's function. It runs the
// command and returns the proper exit code along with
// any error that occurred.
type CommandFunc func(Flags) (int, error)

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
	if cmd.Func == nil && cmd.CobraFunc == nil {
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
	// defaultFactory.Use(func(rootCmd *cobra.Command) {
	// 	rootCmd.AddCommand(caddyCmdToCobra(cmd))
	// })
	commands[cmd.Name] = cmd
}

var commandNameRegex = regexp.MustCompile(`^[a-z0-9]$|^([a-z0-9]+-?[a-z0-9]*)+[a-z0-9]$`)
