package ctxval

type CommonKeyId int

const (
	loggerCtxKey     CommonKeyId = iota
	requestIdCtxKey  CommonKeyId = iota
	requestNumCtxKey CommonKeyId = iota
)
