// Package dao provides a Database Access Object interface to stored data.
//
// All functions require the following variables to be set in the context:
//
// * Logger: for all context-aware logging.
// * Account ID: for multi-tenancy, unless marked with UNSCOPED word.
//
// Functions marked as UNSCOPED can be safely used from contexts where there is
// exactly zero function arguments coming from an user (e.g. ID was retrieved via
// another DAO call that was scoped).
//
package dao

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
)

var GetAccountDao func(ctx context.Context) (AccountDao, error)

// AccountDao TODO
type AccountDao interface {
	Create(ctx context.Context, pk *models.Account) error
	GetById(ctx context.Context, id int64) (*models.Account, error)
	GetOrCreateByIdentity(ctx context.Context, orgId string, accountNumber string) (*models.Account, error)
	GetByOrgId(ctx context.Context, orgId string) (*models.Account, error)
	GetByAccountNumber(ctx context.Context, number string) (*models.Account, error)
	List(ctx context.Context, limit, offset int64) ([]*models.Account, error)
}

var GetPubkeyDao func(ctx context.Context) (PubkeyDao, error)

// PubkeyDao TODO
type PubkeyDao interface {
	Create(ctx context.Context, pk *models.Pubkey) error
	Update(ctx context.Context, pk *models.Pubkey) error
	GetById(ctx context.Context, id int64) (*models.Pubkey, error)
	List(ctx context.Context, limit, offset int64) ([]*models.Pubkey, error)
	Delete(ctx context.Context, id int64) error
}

var GetPubkeyResourceDao func(ctx context.Context) (PubkeyResourceDao, error)

// PubkeyResourceDao interface for manipulating data of Pubkey deployed into cloud provider.
// Allows to track the key in the cloud provider and manage it.
// All methods are unscoped by Account and thus should not get user input as input directly,
// input should first be validated by PubkeyDao.GetById
type PubkeyResourceDao interface {
	UnscopedGetResourceByProviderType(ctx context.Context, pubkeyId int64, provider models.ProviderType) (*models.PubkeyResource, error)
	UnscopedListByPubkeyId(ctx context.Context, pkId int64) ([]*models.PubkeyResource, error)
	UnscopedCreate(ctx context.Context, pkr *models.PubkeyResource) error
	UnscopedDelete(ctx context.Context, id int64) error
}

var GetReservationDao func(ctx context.Context) (ReservationDao, error)

// ReservationDao represents a reservation, an abstraction of one or more background jobs with
// associated detail information different for different cloud providers (like number of vCPUs,
// instance IDs created etc).
type ReservationDao interface {
	// CreateNoop creates no operation reservation with details in a single transaction.
	CreateNoop(ctx context.Context, reservation *models.NoopReservation) error

	// CreateAWS creates AWS reservation with details in a single transaction.
	CreateAWS(ctx context.Context, reservation *models.AWSReservation) error

	// CreateInstance inserts instance associated to a reservation.
	CreateInstance(ctx context.Context, reservation *models.ReservationInstance) error

	// GetById returns reservation for a particular account.
	GetById(ctx context.Context, id int64) (*models.Reservation, error)

	// List returns reservation for a particular account.
	List(ctx context.Context, limit, offset int64) ([]*models.Reservation, error)

	// ListInstances returns instances associated to a reservation. UNSCOPED.
	ListInstances(ctx context.Context, limit, offset int64) ([]*models.ReservationInstance, error)

	// UpdateStatus sets status field and increment step counter by addSteps. UNSCOPED.
	UpdateStatus(ctx context.Context, id int64, status string, addSteps int32) error

	// UpdateReservationIDForAWS updates AWS reservation id field. UNSCOPED.
	UpdateReservationIDForAWS(ctx context.Context, id int64, awsReservationId string) error

	// FinishWithSuccess sets Success flag. UNSCOPED.
	FinishWithSuccess(ctx context.Context, id int64) error

	// FinishWithError sets Success flag and Error flag. UNSCOPED.
	FinishWithError(ctx context.Context, id int64, errorString string) error

	// Delete deletes a reservation. Only used in tests and background cleanup job. UNSCOPED.
	Delete(ctx context.Context, id int64) error
}
