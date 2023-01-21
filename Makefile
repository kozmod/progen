.PHONT: tools
tools: ## Run tools (vet, gofmt, goimports, tidy, etc.)
	@go version
	gofmt -w .
	goimports -w .
	go vet ./...
	@go mod tidy

.PHONT: test
test: ## Run all tests in project with coverage
	@go test ./... -cover

.PHONT: build
build: ## Build binaries with `-ldflags` (Version)
	go build -ldflags "-s -w -X 'github.com/kozmod/progen/internal.Version=$(shell ./version.sh get)'" -o progen .

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help
help: ## List all make targets with description
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'