package pgx

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func init() {
	dao.GetStatDao = getStatDao
}

type statDao struct{}

func getStatDao(ctx context.Context) dao.StatDao {
	return &statDao{}
}

func (x *statDao) getUsage(ctx context.Context, interval string) ([]*models.UsageStat, error) {
	query := `select provider, 'success' as result, count(provider) as count
	from reservations
	where created_at >= now() - cast($1 as interval)
	  and success = true
	group by provider

	union all

	select provider, 'failure' as result, count(provider) as count
	from reservations
	where created_at >= now() - cast($1 as interval)
	  and success = false
	group by provider

	union all

	select provider, 'pending' as result, count(provider) as count
	from reservations
	where created_at >= now() - cast($1 as interval)
	  and success is null
	group by provider`

	var result []*models.UsageStat
	rows, err := db.Pool.Query(ctx, query, interval)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}

	err = pgxscan.ScanAll(&result, rows)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}

	return result, nil
}

func (x *statDao) Get(ctx context.Context) (*models.Statistics, error) {
	usage24h, err := x.getUsage(ctx, "24 hours")
	if err != nil {
		return nil, fmt.Errorf("get usage error: %w", err)
	}

	usage28d, err := x.getUsage(ctx, "28 days")
	if err != nil {
		return nil, fmt.Errorf("get usage error: %w", err)
	}

	return &models.Statistics{
		Usage24h: usage24h,
		Usage28d: usage28d,
	}, nil
}
