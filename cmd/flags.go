package guardcmd

import "github.com/spf13/pflag"

// Flags wraps a FlagSet so that typed values
// from flags can be easily retrieved.
type Flags struct {
	*pflag.FlagSet
}

// String returns the string representation of the
// flag given by name. It panics if the flag is not
// in the flag set.
func (f Flags) String(name string) string {
	return f.FlagSet.Lookup(name).Value.String()
}
