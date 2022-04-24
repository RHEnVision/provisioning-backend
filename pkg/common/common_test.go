// Copyright Red Hat

package common

import (
	"context"
	"net/http"
	"testing"

	"github.com/onsi/gomega"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

func TestGetAccount(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	req, err := http.NewRequest("GET", "/api/provisioning-api/v0/auth_realms", nil)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	// With no account number set in the request's context
	account, err := GetAccount(req)

	g.Expect(account).To(gomega.Equal(""))
	g.Expect(err.Error()).To(gomega.Equal("cannot find account number"))
	
	// Set mock account number in request's context
	sampleHeaderValue := identity.XRHID{
		Identity: identity.Identity {
			AccountNumber: "0369233",
		},
	}
	ctx := context.WithValue(req.Context(), identity.Key, sampleHeaderValue)
	req = req.WithContext(ctx)

	account, err = GetAccount(req)

	g.Expect(account).To(gomega.Equal("0369233"))
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
}
