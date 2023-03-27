package clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccess(t *testing.T) {
	tests := map[string]struct {
		input string
		app   string
		res   string
		verb  string
	}{
		"simple":            {input: "provisioning:b:c", app: "provisioning", res: "b", verb: "c"},
		"empty":             {input: "", app: "", res: "", verb: ""},
		"bad sep":           {input: "a/b/c", app: "", res: "", verb: ""},
		"too few segments":  {input: "a:b", app: "", res: "", verb: ""},
		"too many segments": {input: "a:b:c:d", app: "", res: "", verb: ""},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			acc := NewAccess(tc.input)
			assert.Equal(t, tc.res, acc.Resource)
			assert.Equal(t, tc.verb, acc.Verb)
		})
	}
}

func TestIsAllowed(t *testing.T) {
	tests := map[string]struct {
		input   AccessList
		app     string
		res     string
		verb    string
		allowed bool
	}{
		"single": {input: AccessList{
			NewAccess("provisioning:b:c"),
		}, app: "provisioning", res: "b", verb: "c", allowed: true},
		"multiple": {input: AccessList{
			NewAccess("provisioning:b:c"),
			NewAccess("provisioning:e:f"),
		}, app: "provisioning", res: "e", verb: "f", allowed: true},
		"wildcard resource": {input: AccessList{
			NewAccess("provisioning:*:c"),
		}, app: "provisioning", res: "b", verb: "c", allowed: true},
		"wildcard verb": {input: AccessList{
			NewAccess("provisioning:b:*"),
		}, app: "provisioning", res: "b", verb: "c", allowed: true},
		"wildcard both": {input: AccessList{
			NewAccess("provisioning:*:*"),
		}, app: "provisioning", res: "b", verb: "c", allowed: true},
		"empty": {input: AccessList{}, app: "provisioning", res: "b", verb: "c", allowed: false},
		"not allowed": {input: AccessList{
			NewAccess("a:b:c"),
		}, app: "d", res: "e", verb: "f", allowed: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.IsAllowed(tc.res, tc.verb)
			assert.Equal(t, tc.allowed, got)
		})
	}
}
