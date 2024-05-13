//go:build !test

package pgx

import "github.com/RHEnVision/provisioning-backend/internal/dao"

func init() {
	dao.GetAccountDao = getAccountDao
	dao.GetPubkeyDao = getPubkeyDao
	dao.GetReservationDao = getReservationDao
	dao.GetServiceDao = getServiceDao
	dao.GetStatDao = getStatDao
}
