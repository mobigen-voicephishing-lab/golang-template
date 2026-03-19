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
GOFLAGS ?= $(GOFLAGS:)
BUILD_TIME := $(shell date -u +%Y%m%d.%H%M%S)
# 커밋이 없을 때도 Makefile이 로드되도록 폴백
BUILD_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "nogit")

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
GOOS ?= $(DETECTED_OS)
GOARCH ?= $(DETECTED_ARCH)
BUILD_PLATFORM := $(GOOS)/$(GOARCH)

################################################################################

MODULE_NAME := $(shell head -1 go.mod | awk '{print $$2}')
LDFLAGS = -X '$(MODULE_NAME)/internal/infrastructure/config.Name=$(TARGET)'
LDFLAGS += -X '$(MODULE_NAME)/internal/infrastructure/config.Version=$(VERSION)'
LDFLAGS += -X '$(MODULE_NAME)/internal/infrastructure/config.BuildHash=$(BUILD_HASH)'

################################################################################
##                             Docker PARAMS                                 ##
################################################################################

## Docker Build Versions
DOCKER_BUILD_IMAGE = golang:1.26.1-alpine3.23
DOCKER_BASE_IMAGE = alpine:3.23.3

# Binaries.
TOOLS_BIN_DIR := $(abspath bin)
GO_INSTALL = ./scripts/go_install.sh

MOCKGEN_VER := v0.6.0
MOCKGEN_BIN := mockgen
MOCKGEN := $(TOOLS_BIN_DIR)/$(MOCKGEN_BIN)-$(MOCKGEN_VER)

GOLANGCI_LINT_VER := v2.11.3
GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/$(GOLANGCI_LINT_BIN)-$(GOLANGCI_LINT_VER)

SWAG_VER := v2.0.0-rc5
SWAG_BIN := swag
SWAG := $(TOOLS_BIN_DIR)/$(SWAG_BIN)-$(SWAG_VER)

export GO111MODULE=on

## Checks the code style, tests, builds and bundles.
all: lint build test coverage

.PHONY: lint
lint: $(GOLANGCI_LINT)
	@echo Running golangci-lint
	$(GOLANGCI_LINT) run ./...
	@echo lint success

.PHONY: build build-darwin-arm64 build-linux-amd64
build: ## Build binary for current platform (detected: $(BUILD_PLATFORM))
	@echo Building $(TARGET) for $(BUILD_PLATFORM)
	CGO_ENABLED=1 $(GO) build -ldflags "$(LDFLAGS)" -gcflags all=-trimpath=$(PWD) -asmflags all=-trimpath=$(PWD) \
	     -a -installsuffix cgo -o build/bin/$(TARGET) ./cmd/server

build-darwin-arm64: ## Build for Mac (Apple Silicon, arm64)
	@echo Building $(TARGET) for darwin/arm64
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 $(GO) build -ldflags "$(LDFLAGS)" -gcflags all=-trimpath=$(PWD) -asmflags all=-trimpath=$(PWD) \
	     -a -installsuffix cgo -o build/bin/$(TARGET)-darwin-arm64 ./cmd/server

build-linux-amd64: ## Build for Linux amd64 (e.g. deployment server)
	@echo Building $(TARGET) for linux/amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 $(GO) build -ldflags "$(LDFLAGS)" -gcflags all=-trimpath=$(PWD) -asmflags all=-trimpath=$(PWD) \
	     -a -installsuffix cgo -o build/bin/$(TARGET)-linux-amd64 ./cmd/server

# Generate mocks from the interfaces.
.PHONY: mocks
mocks:  $(MOCKGEN)
	go generate ./...

.PHONY: verify-mocks
verify-mocks:  $(MOCKGEN) mocks
	@if !(git diff --quiet HEAD); then \
		echo "generated files are out of date, run make mocks"; exit 1; \
	fi

# Generate swagger docs from annotations.
.PHONY: swag
swag: $(SWAG)
	@echo Running swag init
	$(SWAG) init --generalInfo cmd/server/main.go --output docs/swagger --parseDependency --parseInternal
	@echo swag success

.PHONY: test
test:
	$(GO) test ./... -v -covermode=count -coverprofile=build/coverage.out

.PHONY: coverage
coverage: test
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
	rm -rf build/bin build/coverage.out build/cov-out.html build/test-output.txt

## Run(Dev)



## --------------------------------------
## Tooling Binaries
## --------------------------------------

$(MOCKGEN): ## Build mockgen.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) go.uber.org/mock/mockgen $(MOCKGEN_BIN) $(MOCKGEN_VER)

$(GOLANGCI_LINT): ## Build golangci-lint.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) github.com/golangci/golangci-lint/v2/cmd/golangci-lint $(GOLANGCI_LINT_BIN) $(GOLANGCI_LINT_VER)

$(SWAG): ## Build swag.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) github.com/swaggo/swag/v2/cmd/swag $(SWAG_BIN) $(SWAG_VER)
