package middleware

import (
	"errors"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	ucontext "github.com/Unleash/unleash-client-go/v3/context"
	"github.com/rs/zerolog/log"
)

func AccountMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rhId := identity.Identity(r.Context())
		orgID := rhId.Identity.OrgID
		accountNumber := rhId.Identity.AccountNumber
		logger := log.Ctx(r.Context()).With().Str("account_number", accountNumber).Str("org_id", orgID).Logger()

		cachedAccount := &models.Account{}
		err := cache.Find(r.Context(), orgID+accountNumber, cachedAccount)
		if errors.Is(err, cache.ErrNotFound) {
			// account not found in cache
			accDao := dao.GetAccountDao(r.Context())

			cachedAccount, err = accDao.GetOrCreateByIdentity(r.Context(), orgID, accountNumber)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to fetch account")
				http.Error(w, err.Error(), 500)
				return
			}

			err = cache.Set(r.Context(), orgID+accountNumber, cachedAccount)
			if err != nil {
				logger.Error().Err(err).Msg("Unable to store account to cache")
				http.Error(w, err.Error(), 500)
				return
			}
		} else if err != nil {
			logger.Error().Err(err).Msg("Cache returned error")
			http.Error(w, err.Error(), 500)
			return
		} else {
			logger.Trace().Int64("account", cachedAccount.ID).Msg("Account cache hit")
		}

		// set contexts - account id
		ctx := identity.WithAccountId(r.Context(), cachedAccount.ID)

		// logger
		newLogger := logger.With().
			Int64("account_id", cachedAccount.ID).
			Logger()
		ctx = newLogger.WithContext(ctx)

		// unleash context
		uctx := ucontext.Context{
			UserId:        cachedAccount.OrgID,
			RemoteAddress: r.RemoteAddr,
			Environment:   config.Unleash.Environment,
			AppName:       version.UnleashAppName,
			Properties:    nil,
		}
		ctx = config.WithUnleashContext(ctx, uctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
