package payloads

import (
	"reflect"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	httpClients "github.com/RHEnVision/provisioning-backend/internal/clients/http"
)

func TestFindUserPayload(t *testing.T) {
	type test struct {
		input error
		want  *userPayload
	}

	tests := []test{
		{
			clients.HttpClientErr,
			&userPayload{500, "unknown backend client error"},
		},
		{
			httpClients.CloneNotFoundErr,
			&userPayload{404, "image builder could not find compose clone"},
		},
		{
			clients.MissingProvisioningSources,
			&userPayload{500, "backend service missing provisioning source"},
		},
	}

	for _, tc := range tests {
		got := findUserPayload(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
