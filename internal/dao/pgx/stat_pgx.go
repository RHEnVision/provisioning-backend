package pgx

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type statDao struct{}

func getStatDao(ctx context.Context) dao.StatDao {
	return &statDao{}
}

func (x *statDao) getUsage(ctx context.Context, iStart, iEnd string) ([]*models.UsageStat, error) {
	query := `select provider, 'success' as result, count(provider) as count
	from reservations
	where created_at between now() - cast($1 as interval) and now() - cast($2 as interval)
	  and success = true
	group by provider

	union all

	select provider, 'failure' as result, count(provider) as count
	from reservations
	where created_at between now() - cast($1 as interval) and now() - cast($2 as interval)
	  and success = false
	group by provider

	union all

	select provider, 'pending' as result, count(provider) as count
	from reservations
	where created_at between now() - cast($1 as interval) and now() - cast($2 as interval)
	  and success is null
	group by provider`

	var result []*models.UsageStat
	rows, err := db.Pool.Query(ctx, query, iStart, iEnd)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}

	err = pgxscan.ScanAll(&result, rows)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}

	return result, nil
}

func (x *statDao) Get(ctx context.Context, delayMin int) (*models.Statistics, error) {
	delay := fmt.Sprintf("%d minutes", delayMin)
	usage24h, err := x.getUsage(ctx, "24 hours "+delay, delay)
	if err != nil {
		return nil, fmt.Errorf("get usage error: %w", err)
	}

	usage28d, err := x.getUsage(ctx, "28 days "+delay, delay)
	if err != nil {
		return nil, fmt.Errorf("get usage error: %w", err)
	}

	return &models.Statistics{
		Usage24h: usage24h,
		Usage28d: usage28d,
	}, nil
}
