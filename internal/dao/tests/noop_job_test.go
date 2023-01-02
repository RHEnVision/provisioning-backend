//go:build integration
// +build integration

package tests

import (
	"testing"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReservationNoopOneSuccess(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)

	type test struct {
		name   string
		action func() *models.NoopReservation
		checks func(*models.Reservation)
	}

	tests := []test{
		{
			name: "SingleStepSuccess",
			action: func() *models.NoopReservation {
				// create new reservation
				res := &models.NoopReservation{
					Reservation: models.Reservation{
						AccountID:  ctxval.AccountId(ctx),
						Steps:      1,
						StepTitles: []string{"Test step"},
						Provider:   models.ProviderTypeNoop,
						Status:     "Created",
					},
				}
				err := reservationDao.CreateNoop(ctx, res)
				require.NoError(t, err)
				require.NotZero(t, res.ID)

				args := jobs.NoopJobArgs{
					AccountID:     ctxval.AccountId(ctx),
					ReservationID: res.ID,
					Fail:          false,
				}
				task := jobs.NoopJobTask.WithArgs(ctx, args)
				err = queue.JobQueue.Add(task)
				require.NoError(t, err)

				return res
			},
			checks: func(updatedRes *models.Reservation) {
				require.True(t, updatedRes.Success.Valid)
				require.True(t, updatedRes.Success.Bool)
				require.Empty(t, updatedRes.Error)
				require.Equal(t, int32(1), updatedRes.Step)
				require.Equal(t, int32(1), updatedRes.Steps)
				require.Equal(t, "No operation finished", updatedRes.Status)
			},
		},
		{
			name: "SingleStepError",
			action: func() *models.NoopReservation {
				// create new reservation
				res := &models.NoopReservation{
					Reservation: models.Reservation{
						AccountID:  ctxval.AccountId(ctx),
						Steps:      1,
						StepTitles: []string{"Test step"},
						Provider:   models.ProviderTypeNoop,
						Status:     "Created",
					},
				}
				err := reservationDao.CreateNoop(ctx, res)
				require.NoError(t, err)
				require.NotZero(t, res.ID)

				args := jobs.NoopJobArgs{
					AccountID:     ctxval.AccountId(ctx),
					ReservationID: res.ID,
					Fail:          true,
				}
				task := jobs.NoopJobTask.WithArgs(ctx, args)
				err = queue.JobQueue.Add(task)
				require.NoError(t, err)

				return res
			},
			checks: func(updatedRes *models.Reservation) {
				require.True(t, updatedRes.Success.Valid)
				require.False(t, updatedRes.Success.Bool)
				require.Equal(t, "job failed on request", updatedRes.Error)
				require.Equal(t, int32(1), updatedRes.Step)
				require.Equal(t, int32(1), updatedRes.Steps)
				require.Equal(t, "No operation finished", updatedRes.Status)
			},
		},
		{
			name: "TwoStepsSuccess",
			action: func() *models.NoopReservation {
				res := &models.NoopReservation{
					Reservation: models.Reservation{
						AccountID:  ctxval.AccountId(ctx),
						Steps:      2,
						StepTitles: []string{"Step one", "Step two"},
						Provider:   models.ProviderTypeNoop,
						Status:     "Created",
					},
				}
				err := reservationDao.CreateNoop(ctx, res)
				require.NoError(t, err)
				require.NotZero(t, res.ID)

				args := jobs.NoopTwoArgs{
					AccountID:     ctxval.AccountId(ctx),
					ReservationID: res.ID,
					Fail1:         false,
					Fail2:         false,
				}
				task := jobs.NoopTwoTask.WithArgs(ctx, args)
				err = queue.JobQueue.Add(task)
				require.NoError(t, err)

				return res
			},
			checks: func(updatedRes *models.Reservation) {
				require.True(t, updatedRes.Success.Valid)
				require.True(t, updatedRes.Success.Bool)
				require.Empty(t, updatedRes.Error)
				require.Equal(t, int32(2), updatedRes.Step)
				require.Equal(t, int32(2), updatedRes.Steps)
				require.Equal(t, "Operation two finished", updatedRes.Status)
			},
		},
		{
			name: "TwoStepsFirstFail",
			action: func() *models.NoopReservation {
				res := &models.NoopReservation{
					Reservation: models.Reservation{
						AccountID:  ctxval.AccountId(ctx),
						Steps:      2,
						StepTitles: []string{"Step one", "Step two"},
						Provider:   models.ProviderTypeNoop,
						Status:     "Created",
					},
				}
				err := reservationDao.CreateNoop(ctx, res)
				require.NoError(t, err)
				require.NotZero(t, res.ID)

				args := jobs.NoopTwoArgs{
					AccountID:     ctxval.AccountId(ctx),
					ReservationID: res.ID,
					Fail1:         true,
					Fail2:         false,
				}
				task := jobs.NoopTwoTask.WithArgs(ctx, args)
				err = queue.JobQueue.Add(task)
				require.NoError(t, err)

				return res
			},
			checks: func(updatedRes *models.Reservation) {
				require.True(t, updatedRes.Success.Valid)
				require.False(t, updatedRes.Success.Bool)
				require.Equal(t, "job failed on request", updatedRes.Error)
				require.Equal(t, int32(1), updatedRes.Step)
				require.Equal(t, int32(2), updatedRes.Steps)
				require.Equal(t, "Operation one finished", updatedRes.Status)
			},
		},
		{
			name: "TwoStepsSecondFail",
			action: func() *models.NoopReservation {
				res := &models.NoopReservation{
					Reservation: models.Reservation{
						AccountID:  ctxval.AccountId(ctx),
						Steps:      2,
						StepTitles: []string{"Step one", "Step two"},
						Provider:   models.ProviderTypeNoop,
						Status:     "Created",
					},
				}
				err := reservationDao.CreateNoop(ctx, res)
				require.NoError(t, err)
				require.NotZero(t, res.ID)

				args := jobs.NoopTwoArgs{
					AccountID:     ctxval.AccountId(ctx),
					ReservationID: res.ID,
					Fail1:         false,
					Fail2:         true,
				}
				task := jobs.NoopTwoTask.WithArgs(ctx, args)
				err = queue.JobQueue.Add(task)
				require.NoError(t, err)

				return res
			},
			checks: func(updatedRes *models.Reservation) {
				require.True(t, updatedRes.Success.Valid)
				require.False(t, updatedRes.Success.Bool)
				require.Equal(t, "job failed on request", updatedRes.Error)
				require.Equal(t, int32(2), updatedRes.Step)
				require.Equal(t, int32(2), updatedRes.Steps)
				require.Equal(t, "Operation two finished", updatedRes.Status)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// create reservation and enqueue job(s)
			res := test.action()

			// read reservation until it is finished (max. 2 seconds)
			var updatedRes *models.Reservation
			var err error
			for i := 0; i < 20; i++ {
				time.Sleep(100 * time.Millisecond)
				updatedRes, err = reservationDao.GetById(ctx, res.ID)
				require.NoError(t, err)
				assert.Equal(t, res.ID, updatedRes.ID)

				if updatedRes.Success.Valid {
					break
				}
			}

			// perform checks
			test.checks(updatedRes)
		})
	}
}
