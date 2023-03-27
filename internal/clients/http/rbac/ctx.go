package rbac

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

type cxtKeyId int

const (
	aclCtxKey cxtKeyId = iota
)

// RbacAcl returns ACL interface, when no ACL is present it returns a list that always evaluates to false.
func Acl(ctx context.Context) clients.RbacAcl {
	value := ctx.Value(aclCtxKey)
	if value == nil {
		return clients.NoPermissionsRbacAcl
	}
	return value.(clients.RbacAcl)
}

// WithAcl returns context copy with ACL interface.
func WithAcl(ctx context.Context, id clients.RbacAcl) context.Context {
	return context.WithValue(ctx, aclCtxKey, id)
}
