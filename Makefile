.PHONT: tools
tools:
	@go version
	go vet ./...
	gofmt -w .
	goimports -w .
	go mod tidy

.PHONT: test
test:
	@go test ./... -cover