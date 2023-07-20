//go:build integration
// +build integration

package tests

import (
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStats(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		res := newNoopReservation()
		err := reservationDao.CreateNoop(ctx, res)
		require.NoError(t, err)

		statDao := dao.GetStatDao(ctx)
		stats, err := statDao.Get(ctx, 0)

		require.NoError(t, err)
		assert.Equal(t, int64(1), stats.Usage24h[0].Count)
	})
}
