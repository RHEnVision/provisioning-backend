package notifications

import (
	"context"

	"github.com/rs/zerolog"
)

type noopNotificationClient struct{}

var _ NotificationClient = &noopNotificationClient{}

func getNoopNotificationClient(ctx context.Context) NotificationClient {
	zerolog.Ctx(ctx).Debug().Msg("Using noop notification client")
	return &noopNotificationClient{}
}

func (s *noopNotificationClient) SuccessfulLaunch(ctx context.Context, reservationId int64) {
	logger := zerolog.Ctx(ctx)
	logger.Warn().Msg("SuccessfulLaunch not started (Notifications not configured)")
}

func (s *noopNotificationClient) FailedLaunch(ctx context.Context, reservationId int64, jobError error) {
	logger := zerolog.Ctx(ctx)
	logger.Warn().Msg("FailedLaunch not started (Notifications not configured)")
}
