module github.com/RHEnVision/provisioning-backend

go 1.18

// go get github.com/georgysavva/scany/v2@v2.0.0-alpha.2
// go get github.com/jackc/tern/v2@v2.0.0-beta.2
// go get github.com/IBM/pgxpoolprometheus@v1.1.0-beta.1

require (
	cloud.google.com/go/compute v1.10.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute v1.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions v1.0.0
	github.com/IBM/pgxpoolprometheus v1.1.0
	github.com/Unleash/unleash-client-go/v3 v3.7.0
	github.com/aws/aws-sdk-go-v2 v1.16.16
	github.com/aws/aws-sdk-go-v2/config v1.17.8
	github.com/aws/aws-sdk-go-v2/credentials v1.12.21
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.15.20
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.63.1
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.19
	github.com/aws/smithy-go v1.13.3
	github.com/deepmap/oapi-codegen v1.11.0
	github.com/exaring/otelpgx v0.2.0
	github.com/georgysavva/scany/v2 v2.0.0
	github.com/getkin/kin-openapi v0.104.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/render v1.0.2
	github.com/go-logr/zerologr v1.2.2
	github.com/go-openapi/runtime v0.24.1
	github.com/go-playground/mold/v4 v4.2.0
	github.com/go-playground/validator/v10 v10.11.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/ilyakaznacheev/cleanenv v1.4.0
	github.com/jackc/pgx-zerolog v0.0.0-20220923130014-7856b90a65ae
	github.com/jackc/pgx/v5 v5.0.4
	github.com/jackc/tern/v2 v2.0.0-beta.2
	github.com/lzap/cloudwatchwriter2 v1.0.0
	github.com/lzap/dejq v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.13.0
	github.com/redhatinsights/app-common-go v1.6.4
	github.com/redhatinsights/platform-go-middlewares v0.20.0
	github.com/riandyrn/otelchi v0.5.0
	github.com/rs/xid v1.4.0
	github.com/rs/zerolog v1.28.0
	github.com/stretchr/testify v1.8.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.36.1
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/jaeger v1.10.0
	go.opentelemetry.io/otel/sdk v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0
	golang.org/x/crypto v0.1.0
	golang.org/x/exp v0.0.0-20221011100225-3a7871e7c0b2
	google.golang.org/api v0.98.0
	google.golang.org/genproto v0.0.0-20221010155953-15ba04fc1c0e
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.1.4 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.0.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork v1.1.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v0.7.0 // indirect
	github.com/BurntSushi/toml v1.2.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.13.6 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/errors v0.20.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/loads v0.21.2 // indirect
	github.com/go-openapi/spec v0.20.7 // indirect
	github.com/go-openapi/strfmt v0.21.3 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-openapi/validate v0.22.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.5.1 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/invopop/yaml v0.2.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/jackc/puddle/v2 v2.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/twmb/murmur3 v1.1.5 // indirect
	go.mongodb.org/mongo-driver v1.10.3 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/contrib v1.10.0 // indirect
	go.opentelemetry.io/otel/metric v0.32.1 // indirect
	golang.org/x/net v0.1.0 // indirect
	golang.org/x/oauth2 v0.0.0-20221006150949-b44042a4b9c1 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/grpc v1.50.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/oleiade/lane.v1 v1.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
