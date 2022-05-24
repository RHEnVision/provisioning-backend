package services

import (
	"github.com/RHEnVision/provisioning-backend/internal/clients/ec2"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

func CreatePubkey(w http.ResponseWriter, r *http.Request) {
	payload := &payloads.PubkeyRequest{}
	if err := render.Bind(r, payload); err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), err))
		return
	}

	// TODO: extract this into job queue
	err := dao.WithTransaction(r.Context(), db.DB, func(tx dao.Transaction) error {
		pkDao, err := dao.GetPubkeyDao(r.Context(), tx)
		if err != nil {
			return payloads.NewInitializeDAOError(r.Context(), "pubkey DAO", err)
		}

		err = pkDao.Create(r.Context(), payload.Pubkey)
		if err != nil {
			return payloads.NewDAOError(r.Context(), "create pubkey", err)
		}

		pkrDao, err := dao.GetPubkeyResourceDao(r.Context(), tx)
		if err != nil {
			return payloads.NewInitializeDAOError(r.Context(), "pubkey resource DAO", err)
		}
		pkr := models.PubkeyResource{
			PubkeyID: payload.Pubkey.ID,
			Provider: models.ProviderTypeAWS,
		}
		err = pkrDao.Create(r.Context(), &pkr)
		if err != nil {
			return payloads.NewDAOError(r.Context(), "create pubkey resource", err)
		}

		client := ec2.NewEC2Client(r.Context())
		pkr.Handle, err = client.ImportPubkey(payload.Pubkey, pkr.FormattedTag())
		if err != nil {
			return payloads.NewAWSError(r.Context(), "import pubkey", err)
		}
		log.Trace().Msgf("Pubkey imported as '%s' with tag '%s'", pkr.Handle, pkr.FormattedTag())

		err = pkrDao.Update(r.Context(), &pkr)
		if err != nil {
			var e *dao.MismatchAffectedError
			if errors.As(err, &e) {
				return payloads.NewNotFoundError(r.Context(), err)
			} else {
				return payloads.NewDAOError(r.Context(), "update pubkey resource", err)
			}
		}

		return nil
	})
	if err != nil {
		var e *payloads.ResponseError
		if errors.As(err, &e) {
			renderError(w, r, e)
		} else {
			renderError(w, r, payloads.NewUnknownError(r.Context(), err))
		}
		return
	}

	if err := render.Render(w, r, payloads.NewPubkeyResponse(payload.Pubkey)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "pubkey", err))
	}
}

func ListPubkeys(w http.ResponseWriter, r *http.Request) {
	pubkeyDao, err := dao.GetPubkeyDao(r.Context(), db.DB)
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "pubkey DAO", err))
		return
	}

	pubkeys, err := pubkeyDao.List(r.Context(), 100, 0)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "list pubkeys", err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewPubkeyListResponse(pubkeys)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "list pubkeys", err))
		return
	}
}

func GetPubkey(w http.ResponseWriter, r *http.Request) {
	id, err := ParseUint64(r, "ID")
	if err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "ID", err))
		return
	}

	pubkeyDao, err := dao.GetPubkeyDao(r.Context(), db.DB)
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "pubkey DAO", err))
		return
	}

	pubkey, err := pubkeyDao.GetById(r.Context(), id)
	if err != nil {
		var e *dao.NoRowsError
		if errors.As(err, &e) {
			renderError(w, r, payloads.NewNotFoundError(r.Context(), err))
		} else {
			renderError(w, r, payloads.NewDAOError(r.Context(), "get pubkey by id", err))
		}
		return
	}

	if err := render.Render(w, r, payloads.NewPubkeyResponse(pubkey)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "pubkey", err))
	}
}

func DeletePubkey(w http.ResponseWriter, r *http.Request) {
	id, err := ParseUint64(r, "ID")
	if err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "ID", err))
		return
	}

	pubkeyDao, err := dao.GetPubkeyDao(r.Context(), db.DB)
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "pubkey DAO", err))
		return
	}

	err = pubkeyDao.Delete(r.Context(), id)
	if err != nil {
		var e *dao.MismatchAffectedError
		if errors.As(err, &e) {
			renderError(w, r, payloads.NewNotFoundError(r.Context(), err))
		} else {
			renderError(w, r, payloads.NewDAOError(r.Context(), "delete pubkey", err))
		}
		return
	}

	render.NoContent(w, r)
}
