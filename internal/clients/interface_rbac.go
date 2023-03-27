package clients

import "context"

// GetRbacClient returns RBAC interface implementation. There are currently
// two implementations available: HTTP and stub. In case the client could not
// be established, the function logs an error and returns an implementation that
// does not allow any permission.
var GetRbacClient func(ctx context.Context) Rbac

// Rbac interface provides access to the RBAC backend service API. Each action that needs to
// check must provide resource (e.g. pubkey) and action (e.g. write) in order to check permission
// presence for principal that is in the identity headers. Definition of permissions and default
// roles are at https://github.com/RedHatInsights/rbac-config (app named "provisioning").
type Rbac interface {
	// GetPrincipalAccess return an ACL object that can be used to check permissions
	GetPrincipalAccess(ctx context.Context) (RbacAcl, error)

	// Ready returns readiness information
	Ready(ctx context.Context) error
}

// RBAC Access Control List is used to determine if current account can perform
// an operation on a particular resource
type RbacAcl interface {
	// IsAllowed checks if current account can perform "verb" on particular "resource"
	IsAllowed(res, verb string) bool
}

// NoPermissionsRbacAcl is an access list which denies all access. This is used in case there is no ACL in context.
var NoPermissionsRbacAcl RbacAcl = noPermAcl{}

// AllPermissionsRbacAcl is an access list which grants all access. This is used in unit tests.
var AllPermissionsRbacAcl RbacAcl = allPermAcl{}

type noPermAcl struct{}

func (r noPermAcl) IsAllowed(_, _ string) bool {
	return false
}

type allPermAcl struct{}

func (r allPermAcl) IsAllowed(_, _ string) bool {
	return true
}
