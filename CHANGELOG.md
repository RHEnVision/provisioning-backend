<a name="unreleased"></a>
## [Unreleased]

### Chores
- **changelog:** GHA action for commits
- **changelog:** introduced changelog generator


<a name="0.8.0"></a>
## [0.8.0] - 2022-10-10


- Update to app-common-go 1.6.4
- Rebase and go mod tidy
- Rename configs/ to config/
- Regenerate example config
- Heartbeat is time interval now
- Logging config is string instead of int
- Remove viper for cleanenv
- Regenerate types and drop duplicities
- Do not add duplicate entries into instance types
- Remove unnecessary error check
- Add the new instance_types to OpenAPI
- Rename clowder service to 'api'
- Add prometheus pgx metrics
- Launch worker alongside api
- Rename architectures X86_64
- Add otel to EC2 client
- Rename Pubkey Resource DAO functions
- return logger as pointer
- Get value without ptr magic
- Set region on pubkey resource Create
- Drop error from GetXXXDAO functions
- Change to native pgx/v5 driver with scany
- Drop signin region, add logging, simplify
- Refresh EC2 types correctly
- Refresh EC2 zone types correctly
- Refresh types for EC2
- Refresh types for Azure
- Update dependencies
- Improve variable naming in schema generator
- Remove multiple versions of same libraries
- Fetch HEAD^ for DAO tests
- OpenTelemetry with logger exporter and Prometheus
- Implement EC2 types precaching
- Add pubkey id to AWS response
- Try standartize the database name
- Pin linter to 1.47.3
- HMSPROV-285: Reservation detail for AWS ([#235](https://github.com/RHEnVision/provisioning-backend/issues/235))
- Refactor Sources and Cloud provider APIs
- Add HTTP response and app-level caching
- Add GCP doc workflow
- Merge Pubkey and Pubkey Resource DAO
- Add Region To PubkeyResource Model
- Add a make/CI target to check migration changes
- Add a make/CI target to check migration changes


<a name="0.7.0"></a>
## [0.7.0] - 2022-09-15


- Add GCP reservation and launch instance job
- Enable gofumpt ([#230](https://github.com/RHEnVision/provisioning-backend/issues/230))
- Use same secret keys in all environments
- Revert "Rename database to match stage"
- Rename database to match stage
- Use ptr in ec2 client too
- Update golangci-lint to latest
- Package ptr with generics
- HMSPROV-239 - add reservation status endpoint
- Fix Viper Default Values
- Add CPU and MEMORY _REQUESTS parameter defaults
- OpenAPI 3.1 is not Swagger
- Avoid appending when copying slices
- Pubkey resource unscoped by account
- Explain our workflow in AWS further
- Add test for Azure instance types service
- Remove unused old API-generated interface
- Add documentation to the client interfaces
- Unexport client implementation types
- Rename and split client models into separate files
- Split architecture and supported features
- Move HTTP "Doers" to http package
- Move errors to http package
- Improve NotFound errors
- Wrap common errors for all clients too
- Use context-aware logging with client field
- Rename common errors and move AWS-related
- Improve logging of AWS client
- Improve error handling of image builder client
- Improve error handling of sources client
- Add azure instance types endpoint
- Change iqe plugin name
- Changed plugin name

### Reverts
- Rename database to match stage


<a name="0.6.0"></a>
## [0.6.0] - 2022-08-31


- ETag middleware with implementation for OpenAPI
- Split Makefile and create devel docs
- Enable compression in clowder by default
- Transparent HTTP gzip compression
- Transparent HTTP gzip compression
- Unify service test assertions and error handling
- Rename and refactor generation tool
- Refactor instance types for Azure, types for Azure
- Add instance types service test
- Refactor integration tests
- Use require.NoError in all DAO tests
- Add missing APP.COMPRESSION ENV variable
- Move dao tests into internal/ and fix config
- Upgrade to Go 1.18
- Region aware reservation endpoint
- Improve EC2 error handling
- Simplify EC2 client initialization
- HMSPROV-79 Resource tests
- List responses need array schemas
- Fix path in makefile and configs
- Add sts and ec2 interfaces and stubs
- Fix stubs structure
- Fix variadic function bug
- Support for multiple names and JSONB detail
- Rework environmental variables


<a name="0.5.0"></a>
## [0.5.0] - 2022-08-30


- Add region to instance type endpoint
- Update types .http file and update log line
- Update types .http file and update log line
- Add storage info to instance type model ([#191](https://github.com/RHEnVision/provisioning-backend/issues/191))
- Use paginator to list instance types
- Add x86_64 mac architecture
- Add test for creating and filtering new instance types
- Remove limitation on number of instance types
- Change instance types endpoint ([#139](https://github.com/RHEnVision/provisioning-backend/issues/139))
- Add reservation name and poweroff flags
- Filter Instance types
- Configurable HTTP client body logging
- Set AWS credentials from kub secret
- Make log field truncating configurable
- add identity header to ListApplicationTypes
- Change source id from int to string
- Use message instead of field for panic messages ([#178](https://github.com/RHEnVision/provisioning-backend/issues/178))
- Drop account endpoints
- Remove unused Update for pubkey resource
- Use Account from headers in NoopService
- Regenerate clients and delete unused test
- Add Instance Type model
- Remove internal only docs
- Use Request/Response suffixes in OpenAPI schemas
- Remove pbworker binary from the repository
- Add missing pbworker Makefile target
- Add version to HTTP header and logs
- HMSPROV-179: Make pubkey name and body unique per account ([#135](https://github.com/RHEnVision/provisioning-backend/issues/135))
- Add reservation to OpenApi
- Add aws_reservation_service_test
- Job queue stub and test for noop job ([#159](https://github.com/RHEnVision/provisioning-backend/issues/159))
- Adding requestBody to POST /pubkyes
- Update docs for ephemeral deploy
- HMSPROV-79 DAO Account tests
- OpenAPI: Pubkey get and list payload fix
- Add fifth missing arg to createPubkeyResource query ([#160](https://github.com/RHEnVision/provisioning-backend/issues/160))
- Add Source ID
- Generalize testing factories


<a name="0.3.0"></a>
## [0.3.0] - 2022-08-11


- Client proxy support with basic auth for stage
- Add golang lint suggestion
- Use table truncate for each test
- Implement setup/teardown methods in idiomatic Go
- Use GetAWSAmi of image builder client
- Fix IB client ([#153](https://github.com/RHEnVision/provisioning-backend/issues/153))
- Fix AccountNumber type to NullString
- Fix incorrect string truncation in logger ([#152](https://github.com/RHEnVision/provisioning-backend/issues/152))
- Update GHA actions to latest versions
- Improve configuration for integration tests
- Change SourcesClientError to ClientError
- Remove pubkey from reservation endpoint
- Remove pubkey from reservation endpoint
- add image builder service for AMI
- Rename InstanceReservation model
- Add dao test for instance_reservation
- Add launch instance aws job
- Add symlinks for CI scripts
- Move CICD files into .rhcicd dir
- HMSPROV-79 Reservation dao test
- OpenAPI spec owned by QE ([#140](https://github.com/RHEnVision/provisioning-backend/issues/140))
- Use arbitrary path in middleware test ([#137](https://github.com/RHEnVision/provisioning-backend/issues/137))
- Clean up unused methods from Sources API client ([#133](https://github.com/RHEnVision/provisioning-backend/issues/133))
- Define source model ourselves ([#132](https://github.com/RHEnVision/provisioning-backend/issues/132))
- Keep only account id in the context ([#131](https://github.com/RHEnVision/provisioning-backend/issues/131))
- HMSPROV-158: Add validations to pubkeys ([#106](https://github.com/RHEnVision/provisioning-backend/issues/106))
- Fix deploy labes for backend app ([#129](https://github.com/RHEnVision/provisioning-backend/issues/129))
- Update pr_check component
- Push sha tag as well as latest ([#125](https://github.com/RHEnVision/provisioning-backend/issues/125))
- HMSPROV-115: Update makefile, add pr_checks and build_deploy ([#110](https://github.com/RHEnVision/provisioning-backend/issues/110))
- Use EC2 client after assuming role
- Integration Test Makefile and Sample
- Add GH action run for DAO tests ([#123](https://github.com/RHEnVision/provisioning-backend/issues/123))
- Remove sources/ID
- Add data for aws launch instance


<a name="0.2.0"></a>
## [0.2.0] - 2022-07-27


- Test identity DAO independent ([#119](https://github.com/RHEnVision/provisioning-backend/issues/119))
- Remove redundant code from sources stub
- HMSPROV-157: Guard Pubkey DAO by tenant ([#105](https://github.com/RHEnVision/provisioning-backend/issues/105))
- Clean up unused sources test fixtures
- Clean up old sources client approach
- Add SourcesClientV2Stub, remove and fix tests
- Implement GetArn and ListProvisioningSources
- Add sources wrapper client
- Integration Test Makefile and Sample
- Fix type for error As ([#111](https://github.com/RHEnVision/provisioning-backend/issues/111))
- Add Account to test context ([#108](https://github.com/RHEnVision/provisioning-backend/issues/108))
- Middleware creates non-existent accounts ([#100](https://github.com/RHEnVision/provisioning-backend/issues/100))
- Improve pubkey DAO stub ([#102](https://github.com/RHEnVision/provisioning-backend/issues/102))
- Fixes UpdatePubkey ([#103](https://github.com/RHEnVision/provisioning-backend/issues/103))
- HMSPROV-108: Return correct status on not found ([#101](https://github.com/RHEnVision/provisioning-backend/issues/101))
- Use ptr package from AWS smithy project
- Split reservation into detail and noop type
- No operation reservation for job testing
- Encapsulate context getters and setters ([#98](https://github.com/RHEnVision/provisioning-backend/issues/98))
- Test account middleware ([#92](https://github.com/RHEnVision/provisioning-backend/issues/92))
- Use full version for sources ([#97](https://github.com/RHEnVision/provisioning-backend/issues/97))
- Use versioned endpoint for sources ([#96](https://github.com/RHEnVision/provisioning-backend/issues/96))
- Add Account DAO stubs
- HMSPROV-137: Get account from identity ([#90](https://github.com/RHEnVision/provisioning-backend/issues/90))
- Updated source seed script to create DB entries ([#85](https://github.com/RHEnVision/provisioning-backend/issues/85))
- Pass errors by value copy ([#94](https://github.com/RHEnVision/provisioning-backend/issues/94))
- NoRowsError wrap the original error
- Add documentation for AWS IAM Role configuration ([#86](https://github.com/RHEnVision/provisioning-backend/issues/86))
- Move AppTypeId to sources package
- Fetch Application Type from sources
- Fix prefix in http client environment ([#88](https://github.com/RHEnVision/provisioning-backend/issues/88))
- Update dependencies 2022-07
- Build the sources URL from Endpoint
- Add status code for sources client error
- Change to meaningful error
- HMSPROV-125: Add pubkeys to the OpenAPI ([#76](https://github.com/RHEnVision/provisioning-backend/issues/76))
- Get sources service URL
- Add list instance types endpoint to spec
- Enforce identity header for instance types
- Reservations endpoint using job queue library dejq ([#5](https://github.com/RHEnVision/provisioning-backend/issues/5))
- Add arn functionality, instance types payload and service
- update spec doc ([#73](https://github.com/RHEnVision/provisioning-backend/issues/73))
- Make identity script MacOS compatible ([#70](https://github.com/RHEnVision/provisioning-backend/issues/70))
- Add sources under dependencies in clowdapp.yaml
- Add route for openapi spec
- add sources to spec file
- Add error responses to the spec file
- Fix to getPubKeyId
- Fix typo
- add openapi docs
- Change List Sources to consider application type
- Validate spec and clients via Makefile
- Add application type endpoints for sources ([#58](https://github.com/RHEnVision/provisioning-backend/issues/58))
- Make scripts to work from any directory
- compare and validate openapi spec
- Scripts for quick sources backend setup
- Add tests for ListSources and GetSource
- Add clean make target
- Use int64 for PKs in OpenAPI spec
- Correct phony make targed and rename to purgedb
- Use int64 for BIGINT PKs
- Add sources integration an interface and an implementation
- add openapi spec
- Add pubkey create tests ([#36](https://github.com/RHEnVision/provisioning-backend/issues/36))
- Add List sources and Get source


<a name="0.1.0"></a>
## 0.1.0 - 2022-06-15


- Replace gomigrate by tern in tools
- Fix and regenerate clients
- Add make for purging DB
- Replace golang-migrate with tern
- Rename test support package
- Add identities to .http test clients
- Add type to identity header generator
- Add identity header generator
- Rename the test support package
- Add support for X-RH identity
- sources service and client
- Split updating and generating of clients ([#39](https://github.com/RHEnVision/provisioning-backend/issues/39))
- Update deploy ephemeral doc for mac users
- Fix bonfire config path
- Add image-builder client
- Add sts client
- Filter out AWS secrets from logs
- Add CONTRIBUTING.md
- Fix variable naming ([#32](https://github.com/RHEnVision/provisioning-backend/issues/32))
- Add unit test for service
- Move pubkey resource upload to a separate action
- Change golangci version in makefile
- Setup clowder in real config
- Variable shadowing is a lint error now
- Drop production compile tag
- Environment aware config initialization
- Configurable mounting point and port
- Simplify transactions
- Always assume production in Clowder
- Use correct image name
- Register pgx logging properly
- Add unit tests for model
- Enforce goimports to prevent import ordering conflicts
- Adds quay push automation
- Use drop schema instead of drop table
- Pubkey and PubkeyResource models, DAO, API
- Rename clouds/ package to clients/ so we can add sources too
- Add testing block to clowdapp spec
- Nest all routes under named path
- Add deploy to Ephemeral configs and guide
- Create pubkey model and a POC REST API (not OpenAPI yet)
- Use config for Cloudwatch AWS credentials
- Fix camelCase for config, add drop_all seed script
- Accounts table, DAO sqlx implementation, seeding
- Use viper for config managemet
- Fix typos in migrations
- Migration respects configuration now
- Add sqlx database connection and fix migrations
- Remove leaked personal AWS token, update .gitignore
- Support for .env and global config
- Basic routing, monitoring, migrations, makefile
- Basic go app skeleton
- Initial commit


[Unreleased]: https://github.com/RHEnVision/provisioning-backend/compare/0.8.0...HEAD
[0.8.0]: https://github.com/RHEnVision/provisioning-backend/compare/0.7.0...0.8.0
[0.7.0]: https://github.com/RHEnVision/provisioning-backend/compare/0.6.0...0.7.0
[0.6.0]: https://github.com/RHEnVision/provisioning-backend/compare/0.5.0...0.6.0
[0.5.0]: https://github.com/RHEnVision/provisioning-backend/compare/0.3.0...0.5.0
[0.3.0]: https://github.com/RHEnVision/provisioning-backend/compare/0.2.0...0.3.0
[0.2.0]: https://github.com/RHEnVision/provisioning-backend/compare/0.1.0...0.2.0
