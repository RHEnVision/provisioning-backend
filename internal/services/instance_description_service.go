package services

import (
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
)

func ListInstancesDescription(w http.ResponseWriter, r *http.Request) {
	id, err := ParseInt64(r, "ID")
	if err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "unable to parse ID parameter", err))
		return
	}

	rDao := dao.GetReservationDao(r.Context())
	reservation, err := rDao.GetById(r.Context(), id)
	if err != nil {
		renderNotFoundOrDAOError(w, r, err, "get reservation detail")
		return
	}

	switch reservation.Provider {
	case models.ProviderTypeUnknown, models.ProviderTypeNoop:
		if err := render.Render(w, r, payloads.NewReservationResponse(reservation)); err != nil {
			renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render reservation", err))
		}
	case models.ProviderTypeAWS:
		// TODO: move to aws_instance_description service
		reservation, err := rDao.GetAWSById(r.Context(), id)
		if err != nil {
			message := fmt.Sprintf("get AWS reservation with id %d", id)
			renderNotFoundOrDAOError(w, r, err, message)
			return
		}
		instances, err := rDao.ListInstances(r.Context(), id)
		if err != nil {
			message := fmt.Sprintf("get reservation with id id %d", id)
			renderNotFoundOrDAOError(w, r, err, message)
			return
		}

		sourcesClient, err := clients.GetSourcesClient(r.Context())
		if err != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), err))
			return
		}

		authentication, err := sourcesClient.GetAuthentication(r.Context(), reservation.SourceID)
		if err != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), err))
			return
		}

		ec2Client, err := clients.GetEC2Client(r.Context(), authentication, "")
		if err != nil {
			renderError(w, r, payloads.NewAWSError(r.Context(), "unable to get AWS EC2 client", err))
			return
		}
		instancesIDList := make([]string, len(instances))
		for i, instance := range instances {
			instancesIDList[i] = instance.InstanceID
		}

		instancesDescriptionList, err := ec2Client.ListInstancesDescription(r.Context(), instancesIDList)
		if err != nil {
			renderError(w, r, payloads.NewAWSError(r.Context(), "unable to get list of instance description", err))
			return
		}

		if err := render.RenderList(w, r, payloads.NewListInstanceDescriptionResponse(instancesDescriptionList)); err != nil {
			renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render instance description", err))
			return
		}
	case models.ProviderTypeAzure:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "azure instance description is not implemented", ProviderTypeNotImplementedError))
	case models.ProviderTypeGCP:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "gcp instance description is not implemented", ProviderTypeNotImplementedError))
	default:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "provider is not supported", ProviderTypeNotImplementedError))
	}
}
