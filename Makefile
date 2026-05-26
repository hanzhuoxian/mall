
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
                   This option is available when using: make build/build.multiarch
                   Example: make build BINS="iam-apiserver iam-authz-server"
  IMAGES           Backend images to make. Default is all of cmd starting with iam-.
                   This option is available when using: make image/image.multiarch/push/push.multiarch
                   Example: make image.multiarch IMAGES="iam-apiserver iam-authz-server"
  REGISTRY_PREFIX  Docker registry prefix. Default is marmotedu. 
                   Example: make push REGISTRY_PREFIX=ccr.ccs.tencentyun.com/marmotedu VERSION=v1.6.2
  PLATFORMS        The multiple platforms to build. Default is linux_amd64 and linux_arm64.
                   This option is available when using: make build.multiarch/image.multiarch/push.multiarch
                   Example: make image.multiarch IMAGES="iam-apiserver iam-pump" PLATFORMS="linux_amd64 linux_arm64"
  VERSION          The version information compiled into binaries.
                   The default is obtained from gsemver or git.
  V                Set to 1 enable verbose build. Default is 0.
endef

export USAGE_OPTIONS


.PHONY: run
run:
	@$(GO) run -ldflags "$(GO_LDFLAGS)" cmd/mall/main.go

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

## build.multiarch: Build source code for multiple platforms. See option PLATFORMS.
.PHONY: build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch


## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"