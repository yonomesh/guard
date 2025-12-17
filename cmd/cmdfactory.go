package guardcmd

import "github.com/spf13/cobra"

type rootCmdFactory struct {
	constructor func() *cobra.Command
	options     []func(*cobra.Command)
}

func newRootCmdFactory(fn func() *cobra.Command) *rootCmdFactory {
	return &rootCmdFactory{
		constructor: fn,
	}
}

func (f *rootCmdFactory) Apply(fn func(cmd *cobra.Command)) {
	f.options = append(f.options, fn)
}

func (f *rootCmdFactory) Build() *cobra.Command {
	o := f.constructor()
	for _, v := range f.options {
		v(o)
	}
	return o
}
