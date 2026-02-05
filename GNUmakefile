APP_NAME    := txeh
BRANCH_NAME := $(shell git rev-parse --abbrev-ref HEAD)
TAG         ?= $(BRANCH_NAME)
IMAGE_NAME  := $(APP_NAME):$(TAG)

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@echo "  build      Build the Go project"
	@echo "  test       Run tests with race detection"
	@echo "  lint       Run linters"
	@echo "  lint-fix   Run linters with auto-fix"
	@echo "  coverage   Generate coverage report"
	@echo "  security   Run security scans"
	@echo "  verify     Run all checks (lint + test + security)"
	@echo "  image      Build Docker image tagged as '$(APP_NAME):$(TAG)'"
	@echo "  format     Format code"
	@echo "  update     Update dependencies"
	@echo "  clean      Clean build artifacts"
	@echo "  all        Run verify and build"
	@echo "  help       Show this help message"
	@echo ""

.PHONY: default
default: help

.PHONY: build
build:
	go build -C txeh/ -ldflags='-s -w -X github.com/txn2/txeh/txeh/cmd.Version={{.Version}}' -trimpath -v -o dist/$(APP_NAME)

.PHONY: test
test:
	go test -race -v ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./...

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
verify: lint test security
	@echo "All checks passed"

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

.PHONY: update
update:
	go version
	go get -u
	go get -u ./...
	go mod tidy

.PHONY: check
check:
	go vet ./...
	staticcheck ./...
	golangci-lint run ./...

.PHONY: format
format:
	gofumpt -e -l -w .
	gofmt -w -s .
