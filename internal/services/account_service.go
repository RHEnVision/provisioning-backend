package services

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
)

func ListAccounts(w http.ResponseWriter, r *http.Request) {
	accountDao, err := dao.GetAccountDao(r.Context(), db.DB)
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "account DAO", err))
		return
	}

	accounts, err := accountDao.List(r.Context(), 100, 0)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "list accounts", err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewAccountListResponse(accounts)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "list accounts", err))
		return
	}
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	id, err := ParseUint64(r, "ID")
	if err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "ID", err))
		return
	}

	accountDao, err := dao.GetAccountDao(r.Context(), db.DB)
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "account DAO", err))
		return
	}

	account, err := accountDao.GetById(r.Context(), id)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "get account by id", err))
		return
	}

	if err := render.Render(w, r, payloads.NewAccountResponse(account)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "account", err))
	}
}
