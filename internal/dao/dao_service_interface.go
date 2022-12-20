package dao

import (
	"context"
)

var GetServiceDao func(ctx context.Context) ServiceDao

// ServiceDao is used for service operations like migrations. All operations are UNSCOPED.
// See pgx/service_pgx.go for documentation.
type ServiceDao interface {
	RecalculatePubkeyFingerprints(ctx context.Context) (int, error)
}
