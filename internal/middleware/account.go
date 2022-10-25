package middleware

import (
	"errors"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	ucontext "github.com/Unleash/unleash-client-go/v3/context"
)

func AccountMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := ctxval.Logger(r.Context())
		rhId := ctxval.Identity(r.Context())
		orgID := rhId.Identity.OrgID
		accountNumber := rhId.Identity.AccountNumber

		cachedAccount, err := cache.FindAccountId(r.Context(), orgID, accountNumber)
		if errors.Is(err, cache.NotFound) {
			// account not found in cache
			accDao := dao.GetAccountDao(r.Context())

			cachedAccount, err = accDao.GetOrCreateByIdentity(r.Context(), orgID, accountNumber)
			if err != nil {
				logger.Error().Err(err).Msgf("Failed to fetch account by org_id=%s/account=%s", orgID, accountNumber)
				http.Error(w, err.Error(), 500)
				return
			}

			err = cache.SetAccountId(r.Context(), orgID, accountNumber, cachedAccount)
			if err != nil {
				logger.Error().Err(err).Msgf("Unable to store account %s to cache", orgID)
				http.Error(w, err.Error(), 500)
				return
			}
		} else if err != nil {
			logger.Error().Err(err).Msgf("Cache returned error")
			http.Error(w, err.Error(), 500)
			return
		}

		// account found in cache
		logger.Trace().Int64("account", cachedAccount.ID).Msg("Account cache hit")

		// set contexts - account id
		ctx := ctxval.WithAccountId(r.Context(), cachedAccount.ID)

		// logger
		newLogger := logger.With().
			Int64("account_id", cachedAccount.ID).
			Str("org_id", cachedAccount.OrgID).
			Str("account_number", cachedAccount.AccountNumber.String).
			Logger()
		ctx = ctxval.WithLogger(ctx, &newLogger)

		// unleash context
		uctx := ucontext.Context{
			UserId:        cachedAccount.OrgID,
			RemoteAddress: r.RemoteAddr,
			Environment:   config.Unleash.Environment,
			AppName:       version.UnleashAppName,
			Properties:    nil,
		}
		ctx = ctxval.WithUnleashContext(ctx, uctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
