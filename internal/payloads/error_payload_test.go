package payloads

import (
	"reflect"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/usrerr"
	"github.com/go-playground/validator/v10"

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
			usrerr.ErrBadRequest400,
			&userPayload{400, "bad request"},
		},
		{
			httpClients.ErrCloneNotFound,
			&userPayload{404, "image clone not found"},
		},
		{
			clients.ErrMissingProvisioningSources,
			&userPayload{500, "sources backend error"},
		},
		{
			validator.ValidationErrors{},
			nil,
		},
	}

	for _, tc := range tests {
		got := findUserPayload(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
