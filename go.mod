module github.com/RHEnVision/provisioning-backend

go 1.16

require (
	github.com/aws/aws-sdk-go-v2 v1.16.4
	github.com/aws/aws-sdk-go-v2/config v1.8.3
	github.com/aws/aws-sdk-go-v2/credentials v1.12.0
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.15.5
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.37.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.6
	github.com/deepmap/oapi-codegen v1.11.0
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/render v1.0.1
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/jackc/pgconn v1.12.1 // indirect
	github.com/jackc/pgx/v4 v4.10.1
	github.com/jmoiron/sqlx v1.3.1
	github.com/lzap/cloudwatchwriter2 v0.0.0-20220422105429-49017f04c285
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/redhatinsights/app-common-go v1.6.1
	github.com/redhatinsights/platform-go-middlewares v0.17.0
	github.com/rs/xid v1.4.0
	github.com/rs/zerolog v1.26.1
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.7.1
)
