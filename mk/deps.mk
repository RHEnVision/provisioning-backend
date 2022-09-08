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
	go mod tidy

# alias for download-deps
.PHONY: prep
prep: download-deps
