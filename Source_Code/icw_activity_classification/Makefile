APP_NAME := icw_activity_classification
GO ?= go
GOIMPORTS_REVISER ?= goimports-reviser
MAIN := ./cmd/server
PKGS := ./...
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)
COVER_DIR := coverage
COVER_PROFILE := $(COVER_DIR)/coverage.out
COVER_HTML := $(COVER_DIR)/coverage.html

.DEFAULT_GOAL := help

.PHONY: help all verify fmt fmt-check tidy tidy-check deps list env vet test test-race cover cover-open build build-all build-linux-amd64 build-linux-arm64 run clean

help:
	@echo "Usage: make <COMMAND>"
	@echo ""
	@echo "Commands:"
	@echo "  make all                Run fmt, tidy, vet, test, and build"
	@echo "  make verify             Run fmt-check, tidy-check, vet, and test"
	@echo "  make fmt                Format Go files"
	@echo "  make fmt-check          Check Go formatting without modifying files"
	@echo "  make tidy               Run go mod tidy"
	@echo "  make tidy-check         Check go.mod/go.sum are tidy"
	@echo "  make deps               Download Go module dependencies"
	@echo "  make list               List Go packages"
	@echo "  make env                Print Go environment"
	@echo "  make vet                Run go vet"
	@echo "  make test               Run Go tests"
	@echo "  make test-race          Run Go tests with race detector"
	@echo "  make cover              Generate coverage report"
	@echo "  make cover-open         Generate and open coverage report"
	@echo "  make build              Build local binary"
	@echo "  make build-all          Build local and linux binaries"
	@echo "  make build-linux-amd64  Build linux/amd64 binary"
	@echo "  make build-linux-arm64  Build linux/arm64 binary"
	@echo "  make run                Run service locally"
	@echo "  make clean              Remove generated build and coverage artifacts"

all: fmt tidy vet test build

verify: fmt-check tidy-check vet test

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

build:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) $(MAIN)

build-all: build build-linux-amd64 build-linux-arm64

build-linux-amd64:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -o $(BIN)-linux-amd64 $(MAIN)

build-linux-arm64:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO) build -o $(BIN)-linux-arm64 $(MAIN)

run:
	./$(BIN)

clean:
	rm -rf $(BIN_DIR) $(COVER_DIR)
