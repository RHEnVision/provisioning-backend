package middleware

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

func AccountMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := ctxval.Logger(r.Context())
		rhId := ctxval.Identity(r.Context())
		orgID := rhId.Identity.OrgID
		accountNumber := rhId.Identity.AccountNumber
		cacheKey := cache.AccountKey{
			OrgID:         orgID,
			AccountNumber: accountNumber,
		}

		var foundAccount *models.Account
		if cachedAccount, ok := cache.FindAccountId(r.Context(), cacheKey); ok {
			// account found in cache
			foundAccount = cachedAccount
			logger.Trace().Int64("account", foundAccount.ID).Msg("Account cache hit")
		} else {
			// account not found in cache
			accDao, err := dao.GetAccountDao(r.Context())
			if err != nil {
				logger.Error().Err(err).Msg("Failed to initialize connection to fetch Account info")
				http.Error(w, http.StatusText(500), 500)
				return
			}

			account, err := accDao.GetOrCreateByIdentity(r.Context(), orgID, accountNumber)
			if err != nil {
				logger.Error().Err(err).Msgf("Failed to fetch account by org_id=%s/account=%s", orgID, accountNumber)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			cache.SetAccountId(r.Context(), cacheKey, account)
			foundAccount = account
		}

		newLogger := logger.With().
			Int64("account_id", foundAccount.ID).
			Str("org_id", foundAccount.OrgID).
			Str("account_number", foundAccount.AccountNumber.String).
			Logger()
		ctx := ctxval.WithAccountId(r.Context(), foundAccount.ID)
		ctx = ctxval.WithLogger(ctx, &newLogger)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
