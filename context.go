package guard

import "context"

type Context struct {
	context.Context
	moduleInstances map[string]Module
	cfg             *Config
	ancestry        []Module
	cleanupFunces   []func()
	exitFuncs       []func(context.Context)
}
