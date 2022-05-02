package ctxval

import "context"

type CommonKeyId int

const (
	LoggerCtxKey         CommonKeyId = iota
	RequestIdCtxKey      CommonKeyId = iota
	RequestNumCtxKey     CommonKeyId = iota
	SshKeyCtxKey         CommonKeyId = iota
	SshKeyResourceCtxKey CommonKeyId = iota
)

func GetValue[T string | uint64](ctx context.Context, key CommonKeyId) T {
	return ctx.Value(key).(T)
}
