APP_NAME	:= txeh
BRANCH_NAME := $(shell git rev-parse --abbrev-ref HEAD)
TAG         ?= $(BRANCH_NAME)
IMAGE_NAME  := $(APP_NAME):$(TAG)

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@echo "  build     Build the Go project"
	@echo "  test      Run all tests with verbose output"
	@echo "  image     Build Docker image tagged as '$(APP_NAME):$(TAG)'"
	@echo "  all       Run vendor, test, and build targets"
	@echo "  help      Show this help message"
	@echo ""

.PHONY: default
default: help

.PHONY: build
build:
	go build -C txeh/ -ldflags='-s -w -X github.com/txn2/txeh/txeh/cmd.Version={{.Version}}' -trimpath -mod=readonly -v -o dist/$(APP_NAME)

.PHONY: test
test:
	go test -v ./... -count=1 -p=1

.PHONY: image
image:
	@echo "building $(APP_NAME):$(IMAGE_NAME)"
	docker build --progress plain -t $(IMAGE_NAME) .

.PHONY: all
all: test build

.PHONY: clean
clean:
	go clean -cache

.PHONY: update
update:
	go version
	go get -u
	go get -u ./...
	go mod tidy

.PHONY: check
check:
	go vet ./... || true
	staticcheck --version
	staticcheck ./... || true
	golangci-lint --version
	golangci-lint run ./... || true
