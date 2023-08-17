.DEFAULT_GOAL=noop
.DELETE_ON_ERROR:

.PHONY: noop
noop:

CI?=false
ifeq ($(CI),true)
VERBOSE?=true
endif

VERBOSE?=false
ifeq ($(VERBOSE),true)
VERBOSE_FLAG=-v
else
VERBOSE_FLAG=
endif

VERSION?=$(shell (git describe --tags --exact-match 2> /dev/null || git rev-parse HEAD) | sed "s/^v//")
.PHONY: version
version:
	@echo $(VERSION)

GO_MODULE=$(shell go list -m)

GO_BUILD_DIR=build
.PHONY: build
build:
ifneq ($(wildcard ./cmd/*),)
	mkdir -p $(GO_BUILD_DIR)
	go build $(VERBOSE_FLAG) -ldflags="-s -w -X main.version=$(VERSION)" -o $(GO_BUILD_DIR) ./cmd/...
endif

.PHONY: test
test:
	go test $(VERBOSE_FLAG) -fullpath -cover -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out -o=coverage.txt
ifeq ($(VERBOSE),true)
	cat coverage.txt
endif
	go tool cover -html=coverage.out -o=coverage.html

.PHONY: generate
generate::
	go generate $(VERBOSE_FLAG) ./...

.PHONY: lint
lint:
	$(MAKE) golangci-lint
	$(MAKE) lint-rules

# version:
# - tag: vX.Y.Z
# - branch: master
# - latest
GOLANGCI_LINT_VERSION?=v1.54.0
# Installation type:
# - binary
# - source
GOLANGCI_LINT_TYPE?=binary

ifeq ($(GOLANGCI_LINT_TYPE),binary)

GOLANGCI_LINT_DIR=$(shell go env GOPATH)/pkg/golangci-lint/$(GOLANGCI_LINT_VERSION)
GOLANGCI_LINT_BIN=$(GOLANGCI_LINT_DIR)/golangci-lint

$(GOLANGCI_LINT_BIN):
	curl $(VERBOSE_FLAG) -fL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOLANGCI_LINT_DIR) $(GOLANGCI_LINT_VERSION)

.PHONY: install-golangci-lint
install-golangci-lint: $(GOLANGCI_LINT_BIN)

else ifeq ($(GOLANGCI_LINT_TYPE),source)

GOLANGCI_LINT_BIN=go run $(VERBOSE_FLAG) github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

install-golangci-lint:

endif

GOLANGCI_LINT_RUN=$(GOLANGCI_LINT_BIN) $(VERBOSE_FLAG) run
.PHONY: golangci-lint
golangci-lint: install-golangci-lint
ifeq ($(CI),true)
	$(GOLANGCI_LINT_RUN)
else
# Fix errors if possible.
	$(GOLANGCI_LINT_RUN) --fix
endif

.PHONY: golangci-lint-cache-clean
golangci-lint-cache-clean: install-golangci-lint
	$(GOLANGCI_LINT_BIN) cache clean

.PHONY: lint-rules
lint-rules:
	# Disallowed files.
	! find . -name ".DS_Store" | grep "."

	# Mandatory files.
	[ -e .gitignore ]
	[ -e README.md ]
	[ -e LICENSE ]
	[ -e CODEOWNERS ]
	[ -e .github/dependabot.yml ]
	[ -e .github/workflows/ci.yml ]
	[ -e .github/workflows/dependabot_auto_merge.yml ]
	[ -e go.mod ]
	[ -e go.sum ]
	[ -e .golangci.yml ]
	[ -e Makefile ]
	[ -e Makefile-common.mk ]

	# Don't use upper case letter in file and directory name.
	# The convention for separator in name is:
	# - file: "_"
	# - directory in "/cmd": "-"
	# - other directory: shouldn't be separated
	! find . -name "*.go" | grep "[[:upper:]]"

	# Use Go 1.21 in go.mod.
	! grep -n "^go " go.mod | grep -v "go 1.21$$"

.PHONY: mod-update
mod-update:
	go get $(VERBOSE_FLAG) -u all
	$(MAKE) mod-tidy

.PHONY: mod-tidy
mod-tidy:
	go mod tidy $(VERBOSE_FLAG)

.PHONY: git-latest-release
git-latest-release:
	@git tag --list --sort=v:refname --format="%(refname:short) => %(creatordate:short)" | tail -n 1

.PHONY: clean
clean:
	git clean -fdX
	go clean -cache -testcache
	$(MAKE) golangci-lint-cache-clean

ifeq ($(CI),true)

GITHUB_REF?=$(error missing GITHUB_REF)
GITHUB_BRANCH=$(shell echo $(GITHUB_REF) | grep -Po "^refs\/heads/\K.+")
GITHUB_TAG=$(shell echo $(GITHUB_REF) | grep -Po "^refs\/tags/\K.+")

CI_LOG_GROUP_START=@echo "::group::$(1)"
CI_LOG_GROUP_END=@echo "::endgroup::"

.PHONY: ci
ci::
	$(call CI_LOG_GROUP_START,env)
	$(MAKE) ci-env
	$(call CI_LOG_GROUP_END)

	$(call CI_LOG_GROUP_START,build)
	$(MAKE) build
	$(call CI_LOG_GROUP_END)

	$(call CI_LOG_GROUP_START,test)
	$(MAKE) test
	$(call CI_LOG_GROUP_END)

	$(call CI_LOG_GROUP_START,lint)
	$(MAKE) lint
	$(call CI_LOG_GROUP_END)

.PHONY: ci-env
ci-env:
	env

ifneq ($(GITHUB_TAG),)
ci::
	$(call CI_LOG_GROUP_START,tag)
	$(MAKE) ci-tag
	$(call CI_LOG_GROUP_END)

GO_PROXY_MODULE_TAG_INFO_URL=https://proxy.golang.org/$(GO_MODULE)/@v/$(GITHUB_TAG).info
.PHONY: ci-tag
ci-tag:
	curl $(VERBOSE_FLAG) -L --fail-with-body $(GO_PROXY_MODULE_TAG_INFO_URL)
# Print an empty line to separate the output of curl and print the log group properly.
	@echo ""
endif

endif # CI end
