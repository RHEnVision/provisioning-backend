package clients

import "strings"

const (
	permissionDelimiter = ":"
	wildcard            = "*"
)

type AccessList []Access

// Access represents a permission. ResourceDefinitions are ignored.
// Inspired by https://github.com/RedHatInsights/rbac-client-go
type Access struct {
	Resource string `json:"resource"`
	Verb     string `json:"verb"`
}

// NewAccess constructs new Access from a string in the form of
// "application:resource:verb". The string may contain wildcards (*).
func NewAccess(access string) Access {
	s := strings.Split(access, permissionDelimiter)
	a := Access{}
	if len(s) == 3 && s[0] == "provisioning" {
		a.Resource = s[1]
		a.Verb = s[2]
		return a
	}

	return a
}

// IsAllowed returns whether an action against a resource is allowed by an AccessList
// taking wildcards into consideration.
func (l AccessList) IsAllowed(res, verb string) bool {
	for _, a := range l {
		if matchWildcard(a.Resource, res) && matchWildcard(a.Verb, verb) {
			return true
		}
	}
	return false
}

func matchWildcard(s1, s2 string) bool {
	return s1 == s2 || s1 == wildcard
}

func (l AccessList) String() string {
	sb := strings.Builder{}
	for _, a := range l {
		sb.WriteString(a.Resource)
		sb.WriteString(":")
		sb.WriteString(a.Verb)
		sb.WriteString(" ")
	}
	return sb.String()
}
