package notifications

import (
	"context"
)

var GetNotificationClient func(ctx context.Context) NotificationClient = getNoopNotificationClient

type NotificationClient interface {
	SuccessfulLaunch(ctx context.Context, reservationId int64)
	FailedLaunch(ctx context.Context, reservationId int64, jobError error)
}
