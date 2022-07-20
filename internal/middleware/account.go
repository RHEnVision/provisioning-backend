package middleware

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
)

func AccountMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rhId := ctxval.Identity(r.Context())
		logger := ctxval.Logger(r.Context())
		accDao, err := dao.GetAccountDao(r.Context())
		if err != nil {
			logger.Error().Err(err).Msg("Failed to initialize connection to fetch Account info")
			http.Error(w, http.StatusText(500), 500)
			return
		}
		acc, err := accDao.GetOrCreateByIdentity(r.Context(), rhId.Identity.OrgID, rhId.Identity.AccountNumber)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to fetch account by org_id=%s/account=%s", rhId.Identity.OrgID, rhId.Identity.AccountNumber)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		newLogger := logger.With().Str("org_id", acc.OrgID).Str("account_number", *acc.AccountNumber).Logger()
		ctx := ctxval.WithAccount(r.Context(), acc)
		ctx = ctxval.WithLogger(ctx, &newLogger)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
