# Copyright (c) 2026 Mobigen JBLIM. All Rights Reserved.

################################################################################
##                             PROGRAM PARAMS                                 ##
################################################################################

# program name and version info
REPO := repo.iris.tools/template/
TARGET := test
VERSION := v1.0.0
IMAGE ?= $(REPO)$(TARGET):$(VERSION)

################################################################################

GO ?= $(shell command -v go 2> /dev/null)
MACHINE = $(shell uname -m)
GOFLAGS ?= $(GOFLAGS:)
BUILD_TIME := $(shell date -u +%Y%m%d.%H%M%S)
BUILD_HASH := $(shell git rev-parse --short HEAD)

# 현재 환경 감지 (멀티 아키텍처 빌드 분기용)
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)
# uname -m: arm64(aarch64), x86_64(amd64) 등 → GOARCH로 매핑
ifeq ($(UNAME_M),x86_64)
  DETECTED_ARCH := amd64
else ifeq ($(UNAME_M),amd64)
  DETECTED_ARCH := amd64
else ifeq ($(UNAME_M),arm64)
  DETECTED_ARCH := arm64
else ifeq ($(UNAME_M),aarch64)
  DETECTED_ARCH := arm64
else
  DETECTED_ARCH := $(UNAME_M)
endif
ifeq ($(UNAME_S),Darwin)
  DETECTED_OS := darwin
else ifeq ($(UNAME_S),Linux)
  DETECTED_OS := linux
else
  DETECTED_OS := $(shell echo $(UNAME_S) | tr '[:upper:]' '[:lower:]')
endif

################################################################################

MODULE_NAME := $(shell head -1 go.mod | awk '{print $$2}')
LDFLAGS = -X '$(MODULE_NAME)/common/appdata.Name=$(TARGET)'
LDFLAGS += -X '$(MODULE_NAME)/common/appdata.Version=$(VERSION)'
LDFLAGS += -X '$(MODULE_NAME)/common/appdata.BuildHash=$(BUILD_HASH)'

################################################################################
##                             Docker PARAMS                                 ##
################################################################################

## Docker Build Versions
DOCKER_BUILD_IMAGE = golang:1.26.1-alpine3.23
DOCKER_BASE_IMAGE = alpine:3.23.3

# Binaries.
TOOLS_BIN_DIR := $(abspath bin)
GO_INSTALL = ./scripts/go_install.sh

MOCKGEN_VER := v1.6.0
MOCKGEN_BIN := mockgen
MOCKGEN := $(TOOLS_BIN_DIR)/$(MOCKGEN_BIN)-$(MOCKGEN_VER)

GOLANGCI_LINT_VER := v2.11.3
GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/$(GOLANGCI_LINT_BIN)-$(GOLANGCI_LINT_VER)

SWAG_VER := v1.16.6
SWAG_BIN := swag
SWAG := $(TOOLS_BIN_DIR)/$(SWAG_BIN)-$(SWAG_VER)

export GO111MODULE=on

## Checks the code style, tests, builds and bundles.
all: lint build test coverage

## Runs golangci-lint (vet, revive/golint 등 정적 분석). 설치: brew install golangci-lint
## 루트에 go/ 폴더가 있으면 "outside main module" 오류가 나므로 삭제할 것 (이미 .gitignore에 있음)
.PHONY: lint
lint: $(GOLANGCI_LINT)
	@echo Running golangci-lint
	$(GOLANGCI_LINT) run ./...
	@echo lint success

# ## Remove golangci-lint v1 binary from bin/ (use system golangci-lint via brew instead).
# .PHONY: lint-clean
# lint-clean:
# 	@rm -f $(TOOLS_BIN_DIR)/golangci-lint $(TOOLS_BIN_DIR)/golangci-lint-* 2>/dev/null || true
# 	@echo "Removed golangci-lint from $(TOOLS_BIN_DIR). Use: brew install golangci-lint"

.PHONY: build build-darwin-arm64 build-linux-amd64
build: ## Build binary for current platform (detected: $(BUILD_PLATFORM))
	@echo Building $(TARGET) for $(BUILD_PLATFORM)
	CGO_ENABLED=1 $(GO) build -ldflags "$(LDFLAGS)" -gcflags all=-trimpath=$(PWD) -asmflags all=-trimpath=$(PWD) \
	     -a -installsuffix cgo -o build/bin/$(TARGET) main.go

build-darwin-arm64: ## Build for Mac (Apple Silicon, arm64)
	@echo Building $(TARGET) for darwin/arm64
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 $(GO) build -ldflags "$(LDFLAGS)" -gcflags all=-trimpath=$(PWD) -asmflags all=-trimpath=$(PWD) \
	     -a -installsuffix cgo -o build/bin/$(TARGET)-darwin-arm64 main.go

build-linux-amd64: ## Build for Linux amd64 (e.g. deployment server)
	@echo Building $(TARGET) for linux/amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 $(GO) build -ldflags "$(LDFLAGS)" -gcflags all=-trimpath=$(PWD) -asmflags all=-trimpath=$(PWD) \
	     -a -installsuffix cgo -o build/bin/$(TARGET)-linux-amd64 main.go

# Generate mocks from the interfaces.
.PHONY: mocks
mocks:  $(MOCKGEN)
	go generate ./...

# Generate swagger docs from annotations.
.PHONY: swag
swag: $(SWAG)
	@echo Running swag init
	$(SWAG) init --output docs/swagger
	@echo swag success

.PHONY: verify-mocks
verify-mocks:  $(MOCKGEN) mocks
	@if !(git diff --quiet HEAD); then \
		echo "generated files are out of date, run make mocks"; exit 1; \
	fi

.PHONY: test
test:
	$(GO) test ./... -v -covermode=count -coverprofile=build/coverage.out

.PHONY: coverage
coverage: ## Run tests and generate HTML coverage report (cov-out.html)
	$(GO) test ./... -covermode=count -coverprofile=build/coverage.out
	$(GO) tool cover -html=build/coverage.out -o build/cov-out.html
	@echo "Coverage report: cov-out.html"

.PHONY: image
image:  ## Build the docker image for linux/amd64 (맥에서 실행해도 결과 이미지는 Linux amd64용)
	@echo "Building $(TARGET) Docker Image (linux/amd64)"
	docker buildx build --platform linux/amd64 \
	--build-arg DOCKER_BUILD_IMAGE=$(DOCKER_BUILD_IMAGE) \
	--build-arg DOCKER_BASE_IMAGE=$(DOCKER_BASE_IMAGE) \
	. -f build/Dockerfile -t $(IMAGE) \
	--no-cache

## Clean Cache
.PHONY: clean
clean:
	go clean -i -cache -testcache

## --------------------------------------
## Tooling Binaries
## --------------------------------------

$(MOCKGEN): ## Build mockgen.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) github.com/golang/mock/mockgen $(MOCKGEN_BIN) $(MOCKGEN_VER)

$(GOLANGCI_LINT): ## Build golangci-lint.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) github.com/golangci/golangci-lint/v2/cmd/golangci-lint $(GOLANGCI_LINT_BIN) $(GOLANGCI_LINT_VER)

$(SWAG): ## Build swag.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) github.com/swaggo/swag/cmd/swag $(SWAG_BIN) $(SWAG_VER)
