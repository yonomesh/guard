package common

import "context"

type (
	// CancelCauseFunc 是 Go 的 context 包中用于取消操作并传递取消原因的函数类型
	ContextCancelCauseFunc = context.CancelCauseFunc
)

var (
	// 它返回一个新的 Context 和一个取消函数（cancelFunc）
	// 通过调用返回的 cancelFunc 可以触发取消操作，并且设置一个原因
	ContextWithCancelCause = context.WithCancelCause
	// 这个函数可以用来查询 Context 取消的原因。
	ContextCause = context.Cause
)
