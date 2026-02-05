APP_NAME    := txeh
MODULE      := $(shell go list -m)
BRANCH_NAME := $(shell git rev-parse --abbrev-ref HEAD)
TAG         ?= $(BRANCH_NAME)
IMAGE_NAME  := $(APP_NAME):$(TAG)

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@echo "  build            - Build the Go project"
	@echo "  test             - Run unit and integration tests"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-integration - Run integration tests"
	@echo "  test-e2e         - Run E2E tests (requires Docker)"
	@echo "  lint             - Run linters"
	@echo "  lint-fix         - Run linters with auto-fix"
	@echo "  fmt              - Format code"
	@echo "  coverage         - Generate coverage report"
	@echo "  security         - Run security scans"
	@echo "  verify           - Quick verification (unit tests)"
	@echo "  verify-full      - Full verification (all tests)"
	@echo "  deps             - Download and verify dependencies"
	@echo "  update           - Update dependencies"
	@echo "  image            - Build Docker image"
	@echo "  clean            - Clean build artifacts"
	@echo "  dead-code        - Check for unreachable code"
	@echo "  mutate           - Run mutation testing"
	@echo "  all              - Run verify and build"
	@echo "  help             - Show this help message"
	@echo ""

.PHONY: default
default: help

.PHONY: build
build:
	go build -C txeh/ -ldflags='-s -w -X github.com/txn2/txeh/txeh/cmd.Version={{.Version}}' -trimpath -v -o dist/$(APP_NAME)

.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	go test -race -v ./...

.PHONY: test-integration
test-integration:
	go test -race -v -tags=integration ./test/integration/...

.PHONY: test-e2e
test-e2e:
	go test -race -v -tags=e2e ./test/e2e/...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./...

.PHONY: fmt
fmt:
	gofmt -s -w .
	goimports -w -local $(MODULE) .

.PHONY: coverage
coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | grep total

.PHONY: security
security:
	gosec ./...
	govulncheck ./...

.PHONY: verify
verify: fmt lint test-unit security
	go mod verify
	@echo "All checks passed"

.PHONY: verify-full
verify-full: fmt lint test security
	go mod verify
	@echo "All checks (including integration) passed"

.PHONY: deps
deps:
	go mod download
	go mod verify

.PHONY: update
update:
	go version
	go get -u ./...
	go mod tidy

.PHONY: image
image:
	@echo "building $(APP_NAME):$(IMAGE_NAME)"
	docker build --progress plain -t $(IMAGE_NAME) .

.PHONY: all
all: verify build

.PHONY: clean
clean:
	rm -f coverage.out coverage.html
	go clean -cache

.PHONY: docs
docs:
	pip install -r requirements-docs.txt
	mkdocs serve

.PHONY: docs-build
docs-build:
	pip install -r requirements-docs.txt
	mkdocs build --strict

# Public API methods not called within the module are excluded.
# See: Hosts.Reload, Hosts.RemoveByComments, Hosts.HostAddressLookup
DEADCODE_EXCLUDE := Hosts\.Reload|Hosts\.RemoveByComments|Hosts\.HostAddressLookup

.PHONY: dead-code
dead-code:
	@OUTPUT=$$(deadcode ./... 2>&1 | grep -Ev '$(DEADCODE_EXCLUDE)') || true; \
	if [ -n "$$OUTPUT" ]; then \
		echo "Dead code detected:"; \
		echo "$$OUTPUT"; \
		exit 1; \
	fi

.PHONY: mutate
mutate:
	gremlins unleash --threshold 60 ./...

.PHONY: check
check:
	go vet ./...
	staticcheck ./...
	golangci-lint run ./...
