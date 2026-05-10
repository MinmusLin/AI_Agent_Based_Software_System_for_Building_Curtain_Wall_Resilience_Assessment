APP_NAME := icw_common
GO ?= go
GOIMPORTS_REVISER ?= goimports-reviser
PKGS := ./...
COVER_DIR := coverage
COVER_PROFILE := $(COVER_DIR)/coverage.out
COVER_HTML := $(COVER_DIR)/coverage.html

.DEFAULT_GOAL := help

.PHONY: help all verify gen-all gen gen-ts fmt fmt-check tidy tidy-check deps list env vet test test-race cover cover-open clean

help:
	@echo "Usage: make <COMMAND>"
	@echo ""
	@echo "Commands:"
	@echo "  make all          Run gen, fmt, tidy, vet, and test"
	@echo "  make verify       Run fmt-check, tidy-check, vet, and test"
	@echo "  make all          Generate Go and TypeScript code from proto files"
	@echo "  make gen          Generate Go code from proto files"
	@echo "  make gen-ts       Generate TypeScript code from proto files"
	@echo "  make fmt          Format Go files"
	@echo "  make fmt-check    Check Go formatting without modifying files"
	@echo "  make tidy         Run go mod tidy"
	@echo "  make tidy-check   Check go.mod/go.sum are tidy"
	@echo "  make deps         Download Go module dependencies"
	@echo "  make list         List Go packages"
	@echo "  make env          Print Go environment"
	@echo "  make vet          Run go vet"
	@echo "  make test         Run Go tests"
	@echo "  make test-race    Run Go tests with race detector"
	@echo "  make cover        Generate coverage report"
	@echo "  make cover-open   Generate and open coverage report"
	@echo "  make clean        Remove generated coverage artifacts"

all: gen fmt tidy vet test

verify: fmt-check tidy-check vet test

gen-all: gen gen-ts

gen:
	rm -rf gen
	buf generate

gen-ts:
	rm -rf ../icw_core_web/src/gen
	buf generate proto --template buf.gen.ts.yaml --path proto/core/common.proto --path proto/core/api/

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

clean:
	rm -rf $(COVER_DIR)
