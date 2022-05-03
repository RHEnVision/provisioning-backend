package ctxval

import "context"

type CommonKeyId int

const (
	LoggerCtxKey     CommonKeyId = iota
	RequestIdCtxKey  CommonKeyId = iota
	RequestNumCtxKey CommonKeyId = iota
	ResourceCtxKey   CommonKeyId = iota
)

func GetStringValue(ctx context.Context, key CommonKeyId) string {
	return ctx.Value(key).(string)
}

func GetUInt64Value(ctx context.Context, key CommonKeyId) uint64 {
	return ctx.Value(key).(uint64)
}
