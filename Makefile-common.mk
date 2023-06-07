.DEFAULT_GOAL=noop
.DELETE_ON_ERROR:

.PHONY: noop
noop:

CI?=false

ENSURE_COMMAND=@ which $(1) > /dev/null || (echo "Install the '$(1)' command. $(2)"; exit 1)

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
	go build -v -ldflags="-s -w -X main.version=$(VERSION)" -o $(GO_BUILD_DIR) ./cmd/...
endif

.PHONY: test
test:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out -o=coverage.txt
	cat coverage.txt
	go tool cover -html=coverage.out -o=coverage.html

.PHONY: generate
generate::
	go generate -v ./...

.PHONY: lint
lint:
	$(MAKE) golangci-lint
	$(MAKE) lint-rules

# version:
# - tag: vX.Y.Z
# - branch: master
# - latest
GOLANGCI_LINT_VERSION?=v1.53.2
# Installation type:
# - binary
# - source
GOLANGCI_LINT_TYPE?=binary

ifeq ($(GOLANGCI_LINT_TYPE),binary)

GOLANGCI_LINT_DIR=$(shell go env GOPATH)/pkg/golangci-lint/$(GOLANGCI_LINT_VERSION)
GOLANGCI_LINT_BIN=$(GOLANGCI_LINT_DIR)/golangci-lint

$(GOLANGCI_LINT_BIN):
	curl -vfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOLANGCI_LINT_DIR) $(GOLANGCI_LINT_VERSION)

.PHONY: install-golangci-lint
install-golangci-lint: $(GOLANGCI_LINT_BIN)

else ifeq ($(GOLANGCI_LINT_TYPE),source)

GOLANGCI_LINT_BIN=go run -v github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

install-golangci-lint:

endif

GOLANGCI_LINT_RUN=$(GOLANGCI_LINT_BIN) -v run
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

.PHONY: ensure-command-pcregrep
ensure-command-pcregrep:
	$(call ENSURE_COMMAND,pcregrep,)

.PHONY: lint-rules
lint-rules: ensure-command-pcregrep
	# Disallowed files.
	! find . -name ".DS_Store" | pcregrep "."

	# Mandatory files.
	[ -e .gitignore ]
	[ -e README.md ]
	[ -e CODEOWNERS ]
	[ -e .github/workflows/ci.yml ]
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
	! find . -name "*.go" | pcregrep "[[:upper:]]"

	# Don't export type/function/variable/constant in main package/test.
	! pcregrep -rnM --include=".+\.go$$" --exclude=".+_test\.go$$" "^package main\n(.*\n)*(type|func|var|const) [[:upper:]]" .
	! pcregrep -rnM --include=".+\.go$$" --exclude=".+_test\.go$$" "^package main\n(.*\n)*(var|const) \(\n((\t.*)?\n)*\t[[:upper:]]" .
	! pcregrep -rn --include=".+_test\.go$$" "^(type|var|const) [[:upper:]]" .
	! pcregrep -rnM --include=".+_test\.go$$" "^(var|const) \(\n((\t.*)?\n)*\t[[:upper:]]" .
	! pcregrep -rn --include=".+_test\.go$$" "^func [[:upper:]]" . | pcregrep -v ":func (Test.*\(t \*testing\.T\)|Benchmark.*\(b \*testing\.B\)|Example.*\(\)) {"

	# Don't declare a var block inside a function.
	! pcregrep -rn --include=".+\.go$$" "^\t+var \($$" .

	# Use Go 1.20 in go.mod.
	! pcregrep -n "^go " go.mod | pcregrep -v "go 1.20$$"

.PHONY: mod-update
mod-update:
	go get -v -u all
	$(MAKE) mod-tidy

.PHONY: mod-tidy
mod-tidy:
	go mod tidy -v

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

	$(call CI_LOG_GROUP_START,apt)
	$(MAKE) ci-apt
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

CI_APT_PACKAGES:=pcregrep
.PHONY: ci-apt
ci-apt:
	sudo apt update
	sudo apt install $(CI_APT_PACKAGES)

ifneq ($(GITHUB_TAG),)
ci::
	$(call CI_LOG_GROUP_START,tag)
	$(MAKE) ci-tag
	$(call CI_LOG_GROUP_END)

GO_PROXY_MODULE_TAG_INFO_URL=https://proxy.golang.org/$(GO_MODULE)/@v/$(GITHUB_TAG).info
.PHONY: ci-tag
ci-tag:
	curl -vL --fail-with-body $(GO_PROXY_MODULE_TAG_INFO_URL)
# Print an empty line to separate the output of curl and print the log group properly.
	@echo ""
endif

endif # CI end
