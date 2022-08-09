package services

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

func CreatePubkey(w http.ResponseWriter, r *http.Request) {
	payload := &payloads.PubkeyRequest{}
	if err := render.Bind(r, payload); err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), err))
		return
	}

	pkDao := dao.GetPubkeyDao(r.Context())

	err := pkDao.Create(r.Context(), payload.Pubkey)
	if err != nil {
		if db.IsPostgresError(err, db.UniqueConstraintErrorCode) != nil {
			renderError(w, r, payloads.PubkeyAlreadyExistsError(r.Context(), err))
		} else {
			renderError(w, r, payloads.NewDAOError(r.Context(), "create pubkey", err))
		}
		return
	}

	if err := render.Render(w, r, payloads.NewPubkeyResponse(payload.Pubkey)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "pubkey", err))
	}
}

func ListPubkeys(w http.ResponseWriter, r *http.Request) {
	pubkeyDao := dao.GetPubkeyDao(r.Context())

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
	id, err := ParseInt64(r, "ID")
	if err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "ID", err))
		return
	}

	pubkeyDao := dao.GetPubkeyDao(r.Context())

	pubkey, err := pubkeyDao.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, dao.ErrNoRows) {
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
	logger := ctxval.Logger(r.Context())
	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "sources client v2", err))
		return
	}

	id, err := ParseInt64(r, "ID")
	if err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "ID", err))
		return
	}

	pubkeyDao := dao.GetPubkeyDao(r.Context())

	pubkey, err := pubkeyDao.GetById(r.Context(), id)
	if err != nil {
		renderError(w, r, payloads.NewNotFoundError(r.Context(), err))
		return
	}

	resources, err := pubkeyDao.UnscopedListResourcesByPubkeyId(r.Context(), pubkey.ID)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "delete pubkey", err))
		return
	}

	for _, res := range resources {
		if res.Provider == models.ProviderTypeAWS {
			if res.Handle != "" {
				logger.Info().Msgf("Deleting pubkey resource ID %v with handle %s", res.ID, res.Handle)
				authentication, errAuth := sourcesClient.GetAuthentication(r.Context(), res.SourceID)
				if errAuth != nil {
					if errors.Is(err, clients.NotFoundErr) {
						renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch arn from sources: application not found", errAuth, 404))
						return
					}
					renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch arn from sources", errAuth, 500))
					return
				}

				ec2Client, errEc2 := clients.GetCustomerEC2Client(r.Context(), authentication, res.Region)
				if errEc2 != nil {
					renderError(w, r, payloads.NewAWSError(r.Context(), "failed to establish ec2 connection", errEc2))
					return
				}

				errDelete := ec2Client.DeleteSSHKey(r.Context(), res.Handle)
				if errDelete != nil {
					renderError(w, r, payloads.NewAWSError(r.Context(), "can't delete public key", errDelete))
					return
				}
			} else {
				logger.Warn().Msgf("Pubkey resource with empty handle: resource ID %d", res.ID)
			}
		} else {
			renderError(w, r, payloads.NewInvalidRequestError(r.Context(), ProviderTypeNotImplementedError))
		}
	}

	err = pubkeyDao.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, dao.ErrAffectedMismatch) {
			renderError(w, r, payloads.NewNotFoundError(r.Context(), err))
		} else {
			renderError(w, r, payloads.NewDAOError(r.Context(), "delete pubkey", err))
		}
		return
	}

	render.NoContent(w, r)
}
