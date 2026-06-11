
# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
all: tidy format lint cover build

# ==============================================================================
# Build options
ROOT_PACKAGE=github.com/hanzhuoxian/mall
VERSION_PACKAGE = github.com/hanzhuoxian/mall/pkg/version


# ==============================================================================
# Includes
include scripts/make-rules/common.mk
include scripts/make-rules/install.mk
include scripts/make-rules/golang.mk


# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  DEBUG            Whether to generate debug symbols. Default is 0.
  BINS             The binaries to build. Default is all of cmd.
                   This option is available when using: make build/build.multi
                   Example: make build BINS="user-server order-server"
  IMAGES           Backend images to make. Default is all of cmd starting with mall-.
                   This option is available when using: make image/image.multi/push/push.multi
                   Example: make image.multi IMAGES="mall-user-server mall-order-server"
  REGISTRY_PREFIX  Docker registry prefix. Default is hanzhuoxian. 
                   Example: make push REGISTRY_PREFIX=ccr.ccs.tencentyun.com/hanzhuoxian VERSION=v1.6.2
  PLATFORMS        The multiple platforms to build. Default is linux_amd64 and linux_arm64.
                   This option is available when using: make build.multi/image.multi/push.multi
                   Example: make image.multi IMAGES="mall-user-server mall-order-server" PLATFORMS="linux_amd64 linux_arm64"
  VERSION          The version information compiled into binaries.
                   The default is obtained from gsemver or git.
  V                Set to 1 enable verbose build. Default is 0.
endef

export USAGE_OPTIONS


.PHONY: run
run:
	@$(GO) run $(GO_LDFLAGS) cmd/mall-user-server/main.go
	@$(GO) run $(GO_LDFLAGS) cmd/mall-api-server/main.go

## tidy: Run go mod tidy to clean up go.mod and go.sum.
.PHONY: tidy
tidy:
	@$(MAKE) go.tidy

## lint: Run golangci-lint to lint the code.
.PHONY: lint
lint:
	@$(MAKE) go.lint

## format: Format the source code.
.PHONY: format
format:
	@$(MAKE) go.format

## build: Build the project.
.PHONY: build
build:
	@$(MAKE) go.build

## test: Run unit test.
.PHONY: test
test:
	@$(MAKE) go.test

## cover: Run unit test and get test coverage.
.PHONY: cover 
cover:
	@$(MAKE) go.test.cover

## build.multi: Build source code for multiple platforms. See option PLATFORMS.
.PHONY: build.multi
build.multi:
	@$(MAKE) go.build.multi

## proto: Generate Go code from proto files.
.PHONY: proto
proto:
	@protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(shell find proto -name "*.proto")


## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"