##@ Testing

TEST_TAGS?=test

.PHONY: test
test: ## Run unit tests
	go test -tags=$(TEST_TAGS) ./...

.PHONY: integration-test
integration-test: ## Run integration tests (require database)
	# "go test pkg1 pkg2" would run tests in parallel causing database locks
	go test --count=1 -v -tags=integration ./internal/dao/tests
	go test --count=1 -v -tags=integration ./internal/migrations/code
