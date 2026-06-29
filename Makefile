APP_NAME=CotoChallenge
PORT=8080
CMD_PATH=./cmd/api
GOLANGCI_VERSION=v1.64.8

.PHONY: run mocks test coverage lint
 
## run: locally execution
run:
	PORT=$(PORT) go run $(CMD_PATH)/main.go

## mocks: repository and service mock generation
mocks:
	@go run go.uber.org/mock/mockgen@latest \
		-source=internal/repository/sale_repository.go \
		-destination=internal/repository/mocks/mock_sale_repository.go \
		-package=mocks
	@go run go.uber.org/mock/mockgen@latest \
		-source=internal/service/service.go \
		-destination=internal/service/mocks/mock_sale_service.go \
		-package=mocks

## test: run tests
test:
	go test ./... -v
 
## coverage: run tests and reports coverage
coverage:
	@echo "Resolving source packages and running strict unit tests..."
	@$(eval PKGS := $(shell go list ./... | grep -v /mocks | grep -v /repository | grep -v /cmd | tr '\n' ',' | sed 's/,$$//'))
	@go test -coverpkg=$(PKGS) -coverprofile=coverage.out ./...
	@echo "--------------------------------------------------------"
	@echo "Coverage Report (Source Code Only):"
	@echo "--------------------------------------------------------"
	@go tool cover -func=coverage.out 

## lint: run code quality analysis and formatting compliance checks
lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION) run ./...
