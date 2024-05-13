//go:build !test

package testing

import "testing"

func TestBuildTag(t *testing.T) {
	t.Fatal("Execute tests with '-tags test' and/or '-tags integration'")
}
