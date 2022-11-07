##@ Go modules

.PHONY: tidy-deps
tidy-deps: ## Cleanup Go modules
	go mod tidy

.PHONY: download-deps
download-deps: ## Download Go modules
	go mod download

.PHONY: update-deps
update-deps: ## Update Go modules to latest versions
	go get -u all
	go get github.com/jackc/puddle/v2@v2.0.0 # https://github.com/jackc/puddle/issues/26
	go get github.com/jackc/tern/v2@v2.0.0-beta.3 # not released yet
	go mod tidy

# alias for download-deps
.PHONY: prep
prep: download-deps
