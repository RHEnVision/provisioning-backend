//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/lzap/dejq"
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
			name: "SingleStepNoError",
			action: func() *models.NoopReservation {
				// create new reservation
				res := &models.NoopReservation{
					Reservation: models.Reservation{
						AccountID:  1,
						Steps:      1,
						StepTitles: []string{"Test step"},
						Provider:   models.ProviderTypeNoop,
						Status:     "Created",
					},
				}
				err := reservationDao.CreateNoop(ctx, res)
				require.NoError(t, err)
				require.NotZero(t, res.ID)

				// create new job
				successJob := dejq.PendingJob{
					Type: queue.TypeNoop,
					Body: &jobs.NoopJobArgs{
						AccountID:     1,
						ReservationID: res.ID,
						Fail:          false,
					},
				}

				err = queue.GetEnqueuer().Enqueue(context.Background(), successJob)
				require.NoError(t, err)

				return res
			},
			checks: func(updatedRes *models.Reservation) {
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
						AccountID:  1,
						Steps:      1,
						StepTitles: []string{"Test step"},
						Provider:   models.ProviderTypeNoop,
						Status:     "Created",
					},
				}
				err := reservationDao.CreateNoop(ctx, res)
				require.NoError(t, err)
				require.NotZero(t, res.ID)

				// create new job
				successJob := dejq.PendingJob{
					Type: queue.TypeNoop,
					Body: &jobs.NoopJobArgs{
						AccountID:     1,
						ReservationID: res.ID,
						Fail:          true,
					},
				}

				err = queue.GetEnqueuer().Enqueue(context.Background(), successJob)
				require.NoError(t, err)

				return res
			},
			checks: func(updatedRes *models.Reservation) {
				require.False(t, updatedRes.Success.Bool)
				require.Equal(t, "job failed on request", updatedRes.Error)
				require.Equal(t, int32(1), updatedRes.Step)
				require.Equal(t, int32(1), updatedRes.Steps)
				require.Equal(t, "No operation finished", updatedRes.Status)
			},
		},
		{
			name: "ThreeStepsNoError",
			action: func() *models.NoopReservation {
				// create new reservation
				res := &models.NoopReservation{
					Reservation: models.Reservation{
						AccountID:  1,
						Steps:      3,
						StepTitles: []string{"Step one", "Step two", "Step three"},
						Provider:   models.ProviderTypeNoop,
						Status:     "Created",
					},
				}
				err := reservationDao.CreateNoop(ctx, res)
				require.NoError(t, err)
				require.NotZero(t, res.ID)

				// create new job
				successJob := dejq.PendingJob{
					Type: queue.TypeNoop,
					Body: &jobs.NoopJobArgs{
						AccountID:     1,
						ReservationID: res.ID,
						Fail:          false,
					},
				}
				err = queue.GetEnqueuer().Enqueue(context.Background(), successJob, successJob, successJob)
				require.NoError(t, err)

				return res
			},
			checks: func(updatedRes *models.Reservation) {
				require.True(t, updatedRes.Success.Bool)
				require.Empty(t, updatedRes.Error)
				require.Equal(t, int32(3), updatedRes.Step)
				require.Equal(t, int32(3), updatedRes.Steps)
				require.Equal(t, "No operation finished", updatedRes.Status)
			},
		},
		{
			name: "ThreeStepsFirstFails",
			action: func() *models.NoopReservation {
				// create new reservation
				res := &models.NoopReservation{
					Reservation: models.Reservation{
						AccountID:  1,
						Steps:      3,
						StepTitles: []string{"Step one", "Step two", "Step three"},
						Provider:   models.ProviderTypeNoop,
						Status:     "Created",
					},
				}
				err := reservationDao.CreateNoop(ctx, res)
				require.NoError(t, err)
				require.NotZero(t, res.ID)

				// create new job
				successJob := dejq.PendingJob{
					Type: queue.TypeNoop,
					Body: &jobs.NoopJobArgs{
						AccountID:     1,
						ReservationID: res.ID,
						Fail:          false,
					},
				}
				failureJob := dejq.PendingJob{
					Type: queue.TypeNoop,
					Body: &jobs.NoopJobArgs{
						AccountID:     1,
						ReservationID: res.ID,
						Fail:          true,
					},
				}
				err = queue.GetEnqueuer().Enqueue(context.Background(), failureJob, successJob, successJob)
				require.NoError(t, err)

				return res
			},
			checks: func(updatedRes *models.Reservation) {
				require.False(t, updatedRes.Success.Bool)
				require.Equal(t, "step skipped: step skipped: job failed on request", updatedRes.Error)
				require.Equal(t, int32(3), updatedRes.Step)
				require.Equal(t, int32(3), updatedRes.Steps)
				require.Equal(t, "No operation finished", updatedRes.Status)
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
