# go source files, ignore vendor directory
PKGS = $(shell go list ./... | grep -v /vendor)
SRC := main.go
BINARY := deviceManagerAPI

# Temporary files to be used, you can changed it calling `make TMP_DIR=/tmp`
TMP_DIR ?= .tmp
COVERAGE_FILE := $(TMP_DIR)/coverage.txt
COVERAGE_HTML_FILE := $(TMP_DIR)/coverage.html

.DEFAULT_GOAL := help
.PHONY: test
test: ## Run all the tests
	@echo "--> Running tests..."
	@mkdir -p $(dir $(COVERAGE_FILE))
	@go test -covermode=atomic -coverprofile=$(COVERAGE_FILE) -race -failfast -timeout=30s $(PKGS)

.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	@echo "--> Creating HTML coverage report at $(COVERAGE_HTML_FILE)..."
	@mkdir -p $(dir $(COVERAGE_FILE)) $(dir $(COVERAGE_HTML_FILE))
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML_FILE)
	@echo "--> Open HTML coverage report: google-chrome $(COVERAGE_HTML_FILE)"

build: ## Build the app
	@echo "--> Building binary artifact ($(BINARY))..."
	@go build -o $(BINARY) $(SRC)

.PHONY: clean
clean: ## Clean all built artifacts
	@echo "--> Cleaning all built artifacts..."
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML_FILE) $(BINARY)
	@go clean
	@go mod tidy -v

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
