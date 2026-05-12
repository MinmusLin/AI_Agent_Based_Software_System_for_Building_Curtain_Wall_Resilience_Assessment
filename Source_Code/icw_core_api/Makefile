APP_NAME := icw_core_api
GO ?= go
GOIMPORTS_REVISER ?= goimports-reviser
DOCKERHUB_NAMESPACE ?= minmuslin
IMAGE_TAG ?= latest
IMAGE_PLATFORM ?= linux/amd64,linux/arm64
IMAGE := $(DOCKERHUB_NAMESPACE)/$(APP_NAME):$(IMAGE_TAG)
MAIN := ./cmd/server
PKGS := ./...
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)
COVER_DIR := coverage
COVER_PROFILE := $(COVER_DIR)/coverage.out
COVER_HTML := $(COVER_DIR)/coverage.html

.DEFAULT_GOAL := help

.PHONY: help all verify image fmt fmt-check tidy tidy-check deps list env vet test test-race cover cover-open build build-darwin build-linux build-windows run clean

help:
	@echo "Usage: make <COMMAND>"
	@echo ""
	@echo "Commands:"
	@echo "  make all            Run fmt, tidy, vet, test, and build"
	@echo "  make verify         Run fmt-check, tidy-check, vet, and test"
	@echo "  make image          Build Docker image and push to Docker Hub"
	@echo "  make fmt            Format Go files"
	@echo "  make fmt-check      Check Go formatting without modifying files"
	@echo "  make tidy           Run go mod tidy"
	@echo "  make tidy-check     Check go.mod/go.sum are tidy"
	@echo "  make deps           Download Go module dependencies"
	@echo "  make list           List Go packages"
	@echo "  make env            Print Go environment"
	@echo "  make vet            Run go vet"
	@echo "  make test           Run Go tests"
	@echo "  make test-race      Run Go tests with race detector"
	@echo "  make cover          Generate coverage report"
	@echo "  make cover-open     Generate and open coverage report"
	@echo "  make build          Build Darwin, Linux, and Windows binaries"
	@echo "  make build-darwin   Build Darwin arm64 and x86_64 binaries"
	@echo "  make build-linux    Build Linux arm64 and x86_64 binaries"
	@echo "  make build-windows  Build Windows arm64 and x86_64 binaries"
	@echo "  make run            Run service locally"
	@echo "  make clean          Remove generated build and coverage artifacts"

all: fmt tidy vet test build

verify: fmt-check tidy-check vet test

image:
	docker login
	docker buildx build --platform $(IMAGE_PLATFORM) --build-context icw_common=../icw_common -f Dockerfile -t $(IMAGE) --push .

fmt:
	$(GO) fmt $(PKGS)
	$(GOIMPORTS_REVISER) -rm-unused -format -recursive -imports-order std,general,company,project -company-prefixes icw_common,icw_core_api,icw_core_biz,icw_activity_classification,icw_activity_reasoning,icw_activity_summary .

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "Go files need formatting:" && gofmt -l . && exit 1)
	$(GOIMPORTS_REVISER) -list-diff -set-exit-status -recursive -imports-order std,general,company,project -company-prefixes icw_common,icw_core_api,icw_core_biz,icw_activity_classification,icw_activity_reasoning,icw_activity_summary .

tidy:
	$(GO) mod tidy

tidy-check:
	@cp go.mod /tmp/$(APP_NAME).go.mod
	@cp go.sum /tmp/$(APP_NAME).go.sum
	@$(GO) mod tidy
	@cmp -s go.mod /tmp/$(APP_NAME).go.mod || (echo "go.mod is not tidy; run make tidy" && mv /tmp/$(APP_NAME).go.mod go.mod && mv /tmp/$(APP_NAME).go.sum go.sum && exit 1)
	@cmp -s go.sum /tmp/$(APP_NAME).go.sum || (echo "go.sum is not tidy; run make tidy" && mv /tmp/$(APP_NAME).go.mod go.mod && mv /tmp/$(APP_NAME).go.sum go.sum && exit 1)
	@rm -f /tmp/$(APP_NAME).go.mod /tmp/$(APP_NAME).go.sum

deps:
	$(GO) mod download

list:
	$(GO) list $(PKGS)

env:
	$(GO) env

vet:
	$(GO) vet $(PKGS)

test:
	$(GO) test $(PKGS)

test-race:
	$(GO) test -race $(PKGS)

cover:
	mkdir -p $(COVER_DIR)
	$(GO) test -coverprofile=$(COVER_PROFILE) $(PKGS)
	$(GO) tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	$(GO) tool cover -func=$(COVER_PROFILE)

cover-open: cover
	open $(COVER_HTML)

build: build-darwin build-linux build-windows

build-darwin:
	mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GO) build -o $(BIN)_Darwin_arm64 $(MAIN)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GO) build -o $(BIN)_Darwin_x86_64 $(MAIN)

build-linux:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO) build -o $(BIN)_Linux_arm64 $(MAIN)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -o $(BIN)_Linux_x86_64 $(MAIN)

build-windows:
	mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 $(GO) build -o $(BIN)_Windows_arm64.exe $(MAIN)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build -o $(BIN)_Windows_x86_64.exe $(MAIN)

run:
	./$(BIN)

clean:
	rm -rf $(BIN_DIR) $(COVER_DIR)
