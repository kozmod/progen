.PHONT: tools
tools: ## Run tools (vet, gofmt, goimports, tidy, etc.)
	@go version
	go vet ./...
	gofmt -w .
	goimports -w .
	@go mod tidy

.PHONT: test
test: ## Run all tests in project with coverage
	@go test ./... -cover

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help
help: ## List all make targets with description
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'