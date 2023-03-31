// Provides DAO implementation via the native Postgres driver with scany library.
package pgx

import "github.com/RHEnVision/provisioning-backend/internal/telemetry"

const TraceName = telemetry.TracePrefix + "internal/dao/pgx"
