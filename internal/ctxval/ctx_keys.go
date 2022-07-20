package ctxval

type commonKeyId int

const (
	loggerCtxKey     commonKeyId = iota
	requestIdCtxKey  commonKeyId = iota
	requestNumCtxKey commonKeyId = iota
	accountCtxKey    commonKeyId = iota
)
