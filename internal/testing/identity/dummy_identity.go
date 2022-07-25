package identity

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/aws/smithy-go/ptr"
	rhidentity "github.com/redhatinsights/platform-go-middlewares/identity"
)

var xRhId = newIdentity(stubs.DefaultOrgId, ptr.String(stubs.DefaultAccountNumber))

func AddIdentityHeader(t *testing.T, req *http.Request) *http.Request {
	req.Header.Add("X-Rh-Identity", setUpValidIdentity(t))
	return req
}

func WithCustomIdentity(t *testing.T, ctx context.Context, orgId string, accountNumber *string) context.Context {
	return context.WithValue(ctx, rhidentity.Key, newIdentity(orgId, accountNumber))
}

func WithTenant(t *testing.T, ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, rhidentity.Key, xRhId)
	ctx = stubs.WithAccountDaoOne(ctx)
	accDao, err := dao.GetAccountDao(ctx)
	if err != nil {
		t.Errorf("failed to initialize account %v", err)
	}
	acc, err := accDao.GetByOrgId(ctx, stubs.DefaultOrgId)
	if err != nil {
		t.Errorf("failed to fetch account for default identity %v", err)
	}
	return ctxval.WithAccount(ctx, acc)
}

func newIdentity(orgId string, accountNumber *string) rhidentity.XRHID {
	id := rhidentity.XRHID{
		Identity: rhidentity.Identity{
			OrgID: orgId,
		},
	}
	if accountNumber != nil {
		id.Identity.AccountNumber = *accountNumber
	}
	return id
}

// setUpValidIdentity returns a base64 encoded valid identity.
func setUpValidIdentity(t *testing.T) string {
	jsonIdentity, err := json.Marshal(xRhId)
	if err != nil {
		t.Errorf(`could not marshal test identity to JSON: %s`, err)
	}

	base64Identity := base64.StdEncoding.EncodeToString(jsonIdentity)
	return base64Identity
}
