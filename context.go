package uni

import "context"

type Context struct {
	context.Context
	moduleInstances map[string]Module
	cfg             *Config
	ancestry        []Module
	cleanupFunces   []func()
	exitFuncs       []func(context.Context)
}

type eventEmitter interface {
	Emit(ctx Context, eventName string, data map[string]any) Event
}
