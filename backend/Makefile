.PHONY: swag ### generate swagger docs
swag:
	go tool swag init -g 'internal/app/app.go' --parseInternal --parseDependency
	go tool swag fmt

.PHONY: bin-deps ### install binary dependencies
bin-deps: 
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2

.PHONY: lint ### check by golangci linter
lint: 
	golangci-lint run

.PHONY: generate-mocks ### run mock generate
generate-mocks:
	go generate ./...