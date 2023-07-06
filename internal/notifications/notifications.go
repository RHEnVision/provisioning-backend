package notifications

import (
	"context"
	"encoding/json"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/kafka"
	"github.com/rs/zerolog"
)

type client struct{}

func getNotificationClient(ctx context.Context) NotificationClient {
	zerolog.Ctx(ctx).Debug().Msg("Using kafka notification client")
	return &client{}
}

func Initialize(ctx context.Context) {
	if config.Application.Notifications.Enabled {
		zerolog.Ctx(ctx).Debug().Msg("Initialized kafka notification client")
		GetNotificationClient = getNotificationClient
	} else {
		zerolog.Ctx(ctx).Debug().Msg("Initialized noop notification client")
	}
}

func (x *client) SuccessfulLaunch(ctx context.Context, reservationId int64) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("Triggering a successful launch notification")
	rDao := dao.GetReservationDao(ctx)
	reservation, err := rDao.GetById(ctx, reservationId)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to find reservation by id")
		return
	}
	instances, err := rDao.ListInstances(ctx, reservationId)
	if err != nil {
		logger.Warn().Err(err).Msg("Unable to get instances for reservation")
	}
	NotificationInstancesEvents := make([]kafka.NotificationEvent, len(instances))
	for i, instance := range instances {
		marshalInstance, er := json.Marshal(instance)
		if er != nil {
			logger.Error().Err(err).Msg("Unable to marshal instance")
			return
		}
		NotificationInstancesEvents[i] = kafka.NotificationEvent{Payload: marshalInstance}
	}

	notificationMsg, err := kafka.NotificationMessage{
		Context:   kafka.NotificationContext{Provider: reservation.Provider.String(), LaunchID: reservationId},
		EventType: kafka.NotificationSuccessEventType, Events: NotificationInstancesEvents,
	}.GenericMessage(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to create notification message")
		return
	}
	logger.Info().Msgf("Sending notification message")
	err = kafka.Send(ctx, &notificationMsg)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to send notification message via kafka")
	}
}

func (x *client) FailedLaunch(ctx context.Context, reservationId int64, jobError error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("Triggering a failed launch notification")
	rDao := dao.GetReservationDao(ctx)
	reservation, err := rDao.GetById(ctx, reservationId)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to find reservation by id")
		return
	}
	marshalError, err := json.Marshal(kafka.NotificationError{Error: jobError.Error()})
	if err != nil {
		logger.Error().Err(err).Msg("Unable to marshal error")
		return
	}

	notificationEvent := []kafka.NotificationEvent{{Payload: marshalError}}
	notificationMsg, err := kafka.NotificationMessage{Context: reservation, EventType: kafka.NotificationFailureEventType, Events: notificationEvent}.GenericMessage(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to create notification failure message")
		return
	}
	logger.Info().Msg("Sending notification message")
	err = kafka.Send(ctx, &notificationMsg)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to send notification message via kafka")
	}
}
