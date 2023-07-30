package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	httpClients "github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/page"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

var ErrMissingNameOrBody = errors.New("name or body missing")

func CreatePubkey(w http.ResponseWriter, r *http.Request) {
	payload := &payloads.PubkeyRequest{}
	if err := render.Bind(r, payload); err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "create pubkey", err))
		return
	}

	if payload.Name == "" || payload.Body == "" {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), ErrMissingNameOrBody.Error(), ErrMissingNameOrBody))
	}

	pkDao := dao.GetPubkeyDao(r.Context())

	pk := payload.NewModel()
	err := pkDao.Create(r.Context(), pk)
	if err != nil {
		if db.IsPostgresError(err, db.UniqueConstraintErrorCode) != nil {
			renderError(w, r, payloads.PubkeyDuplicateError(r.Context(), "pubkey with such name or fingerprint already exists for this account", err))
		} else {
			renderError(w, r, payloads.NewDAOError(r.Context(), "create pubkey", err))
		}
		return
	}

	if err := render.Render(w, r, payloads.NewPubkeyResponse(pk)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render pubkey", err))
	}
}

func ListPubkeys(w http.ResponseWriter, r *http.Request) {
	pubkeyDao := dao.GetPubkeyDao(r.Context())

	offset := page.Offset(r.Context()).Int64()
	limit := page.Limit(r.Context()).Int64()

	pubkeys, err := pubkeyDao.List(r.Context(), limit, offset)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "list pubkeys", err))
		return
	}

	if err := render.Render(w, r, payloads.NewPubkeyListResponse(pubkeys)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render pubkeys list", err))
		return
	}
}

func GetPubkey(w http.ResponseWriter, r *http.Request) {
	id, err := ParseInt64(r, "ID")
	if err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "unable to parse ID parameter", err))
		return
	}

	pubkeyDao := dao.GetPubkeyDao(r.Context())

	pubkey, err := pubkeyDao.GetById(r.Context(), id)
	if err != nil {
		message := fmt.Sprintf("get pubkey with id %d", id)
		renderNotFoundOrDAOError(w, r, err, message)
		return
	}

	if err := render.Render(w, r, payloads.NewPubkeyResponse(pubkey)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render pubkey", err))
	}
}

func DeletePubkey(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())
	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	id, err := ParseInt64(r, "ID")
	if err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "unable to parse ID parameter", err))
		return
	}

	pubkeyDao := dao.GetPubkeyDao(r.Context())

	pubkey, err := pubkeyDao.GetById(r.Context(), id)
	if err != nil {
		message := fmt.Sprintf("get pubkey with id %d", id)
		renderNotFoundOrDAOError(w, r, err, message)
		return
	}

	resources, err := pubkeyDao.UnscopedListResourcesByPubkeyId(r.Context(), pubkey.ID)
	if err != nil {
		message := fmt.Sprintf("list resources by pubkey id %d", pubkey.ID)
		renderNotFoundOrDAOError(w, r, err, message)
		return
	}

	for _, res := range resources {
		if res.Provider == models.ProviderTypeAWS {
			if res.Handle != "" {
				logger.Info().Msgf("Deleting pubkey resource ID %v with handle %s", res.ID, res.Handle)
				authentication, errAuth := sourcesClient.GetAuthentication(r.Context(), res.SourceID)
				if errAuth == nil {
					ec2Client, errEc2 := clients.GetEC2Client(r.Context(), authentication, res.Region)
					if errEc2 != nil {
						renderError(w, r, payloads.NewAWSError(r.Context(), "unable to get AWS client", errEc2))
						return
					}

					errDelete := ec2Client.DeleteSSHKey(r.Context(), res.Handle)
					if errDelete != nil {
						renderError(w, r, payloads.NewAWSError(r.Context(), "unable to delete AWS public key", errDelete))
						return
					}
				} else if errors.Is(errAuth, httpClients.AuthenticationForSourcesNotFoundErr) {
					logger.Warn().Msgf("Skipping source %s authorization which is no longer available", res.SourceID)
				} else {
					logger.Warn().Err(errAuth).Msg("Skipping source authorization because sources returned an error")
				}
			} else {
				logger.Warn().Msgf("Skipping pubkey resource %d with empty handle", res.ID)
			}
		} else {
			renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "delete not implemented for this provider", ProviderTypeNotImplementedError))
		}
	}

	err = pubkeyDao.Delete(r.Context(), id)
	if err != nil {
		message := fmt.Sprintf("pubkey with id %d", id)
		renderNotFoundOrDAOError(w, r, err, message)
		return
	}

	render.NoContent(w, r)
}
