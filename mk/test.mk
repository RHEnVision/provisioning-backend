##@ Testing

TEST_TAGS?=test

.PHONY: test
test: check-go ## Run unit tests
	$(GO) test -tags=$(TEST_TAGS) ./...

.PHONY: integration-test
integration-test: check-go ## Run integration tests (require database)
	# "go test pkg1 pkg2" would run tests in parallel causing database locks
	$(GO) test --count=1 -v -tags=integration ./internal/dao/tests
	$(GO) test --count=1 -v -tags=integration ./internal/migrations/code
	$(GO) test --count=1 -v -tags=integration ./internal/queue/tests
