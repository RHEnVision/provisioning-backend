package ctxval

type CommonKeyId int

const (
	LoggerCtxKey     CommonKeyId = iota
	RequestIdCtxKey  CommonKeyId = iota
	RequestNumCtxKey CommonKeyId = iota
	ResourceCtxKey   CommonKeyId = iota
)
