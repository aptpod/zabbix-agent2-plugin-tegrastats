PROJECT_NAME := zabbix-agent2-plugin-tegrastats
BUILD_DIR := target
BUILD_BIN_DIR := $(BUILD_DIR)/$(shell go env GOOS)_$(shell go env GOARCH)
BUILD_TEST_DIR := $(BUILD_DIR)/test

GO ?= go
GOTESTSUM_VERSION := 1.12.0
GOLINT_VERSION := 1.59.1

.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building $(PROJECT_NAME)..."
	GOPROXY=direct $(GO) build -o $(BUILD_BIN_DIR)/$(PROJECT_NAME) $(filter-out %_test.go, $(wildcard *.go))
	@echo "✅ $(BUILD_BIN_DIR)/$(PROJECT_NAME)"

.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -rf $(BUILD_DIR)

.PHONY: test test-unit test-lint test-fmt test-vuln
test: test-unit test-lint test-fmt test-vuln
test-unit:
	@echo "Running unit tests..."
	$(GO) test -v ./...
	@echo "✅"
test-lint:
	@echo "Running lint..."
	@mkdir -p $(BUILD_TEST_DIR)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BUILD_TEST_DIR) v$(GOLINT_VERSION)
	$(BUILD_TEST_DIR)/golangci-lint run --disable=errcheck --timeout=30m --exclude-files=.*_test.go
	@echo "✅"
test-fmt:
	@echo "Running format check..."
	$(GO)fmt -d -e *.go | tee $(BUILD_TEST_DIR)/fmt.diff;
	@! test -s $(BUILD_TEST_DIR)/fmt.diff
	@echo "✅"
test-vuln:
	@echo "Running vuln check..."
	@$(GO) run golang.org/x/vuln/cmd/govulncheck@latest -version | head -n 5
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@echo "✅"
