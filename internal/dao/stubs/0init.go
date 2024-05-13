//go:build test

package stubs

import "github.com/RHEnVision/provisioning-backend/internal/dao"

func init() {
	dao.GetAccountDao = getAccountDao
	dao.GetPubkeyDao = getPubkeyDao
	dao.GetReservationDao = getReservationDao
}
