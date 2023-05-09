// Provides context value operations.
package ctxval

type commonKeyId int

const (
	loggerCtxKey         commonKeyId = iota
	requestIdCtxKey      commonKeyId = iota
	accountIdCtxKey      commonKeyId = iota
	unleashContextCtxKey commonKeyId = iota
	edgeRequestIdCtxKey  commonKeyId = iota
	correlationCtxKey    commonKeyId = iota
)
