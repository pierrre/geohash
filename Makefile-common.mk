.DEFAULT_GOAL=all
.DELETE_ON_ERROR:

NULL:=
SPACE:=$(NULL) $(NULL)

.PHONY: all
all: build test lint

CI?=false
ifeq ($(CI),true)
VERBOSE?=true
TEST_FULLPATH?=true
TEST_COVER?=true
endif

VERBOSE?=false
ifeq ($(VERBOSE),true)
VERBOSE_FLAG=$(SPACE)-v
else
VERBOSE_FLAG=
endif

RACE?=false
ifeq ($(RACE),true)
RACE_FLAG=$(SPACE)-race
else
RACE_FLAG=
endif

VERSION?=$(shell (git describe --tags --exact-match 2> /dev/null || git rev-parse HEAD) | sed "s/^v//")
.PHONY: version
version:
	@echo $(VERSION)

GO?=go
GO_RUN=$(GO) run$(VERBOSE_FLAG)$(RACE_FLAG)
GO_GET=$(GO) get$(VERBOSE_FLAG)
GO_LIST=$(GO) list$(VERBOSE_FLAG)
GO_MOD=$(GO) mod
GO_TOOL=$(GO) tool
GO_TOOL_COVER=$(GO_TOOL) cover

GO_MODULE=$(shell $(GO_LIST) -m)

GO_TAGS?=
ifneq ($(GO_TAGS),)
GO_TAGS_FLAG=$(SPACE)-tags=$(GO_TAGS)
else
GO_TAGS_FLAG=
endif

BUILD_DIR=build
.PHONY: build
build:
ifneq ($(wildcard ./cmd/*/*.go),)
	mkdir -p $(BUILD_DIR)
	$(GO) build$(VERBOSE_FLAG)$(RACE_FLAG)$(GO_TAGS_FLAG) -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR) ./cmd/...
endif

TEST_FULLPATH?=false
ifeq ($(TEST_FULLPATH),true)
TEST_FULLPATH_FLAG=$(SPACE)-fullpath
else
TEST_FULLPATH_FLAG=
endif
TEST_COVER?=false
ifeq ($(TEST_COVER),true)
TEST_COVER_FLAGS=$(SPACE)-cover -coverprofile=coverage.out
else
TEST_COVER_FLAGS=
endif
TEST_COUNT?=
ifneq ($(TEST_COUNT),)
TEST_COUNT_FLAG=$(SPACE)-count=$(TEST_COUNT)
else
TEST_COUNT_FLAG=
endif
.PHONY: test
test:
	$(GO) test$(VERBOSE_FLAG)$(RACE_FLAG)$(TEST_FULLPATH_FLAG)$(GO_TAGS_FLAG)$(TEST_COVER_FLAGS)$(TEST_COUNT_FLAG) ./...
ifeq ($(TEST_COVER),true)
	$(GO_TOOL_COVER) -func=coverage.out -o=coverage.txt
ifeq ($(VERBOSE),true)
	cat coverage.txt
endif
	$(GO_TOOL_COVER) -html=coverage.out -o=coverage.html
endif

.PHONY: generate
generate::
	$(GO) generate$(VERBOSE_FLAG) ./...

.PHONY: lint
lint:
	$(MAKE) golangci-lint
	$(MAKE) lint-rules
	$(MAKE) mod-tidy

# version:
# - tag: vX.Y.Z
# - branch: master
# - latest
GOLANGCI_LINT_VERSION?=v2.8.0
# Installation type:
# - binary
# - source
GOLANGCI_LINT_TYPE?=binary

ifeq ($(GOLANGCI_LINT_TYPE),binary)

GOLANGCI_LINT_DIR=$(shell $(GO) env GOPATH)/pkg/golangci-lint/$(GOLANGCI_LINT_VERSION)
GOLANGCI_LINT_BIN=$(GOLANGCI_LINT_DIR)/golangci-lint

$(GOLANGCI_LINT_BIN):
	curl$(VERBOSE_FLAG) -fL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh  | sh -s -- -b $(GOLANGCI_LINT_DIR) $(GOLANGCI_LINT_VERSION)

.PHONY: install-golangci-lint
install-golangci-lint: $(GOLANGCI_LINT_BIN)

else ifeq ($(GOLANGCI_LINT_TYPE),source)

GOLANGCI_LINT_BIN=$(GO_RUN) github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

install-golangci-lint:

endif

GOLANGCI_LINT_RUN=$(GOLANGCI_LINT_BIN)$(VERBOSE_FLAG) run
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
CHECK_MISSING_FILE=@[ -e $(1) ] || (echo "$(1) file is missing" && false)
lint-rules:
# Disallowed files.
	@! find . -name ".DS_Store" | (grep "." && echo "Disallowed files")

# Mandatory files.
	$(call CHECK_MISSING_FILE,.gitignore)
	$(call CHECK_MISSING_FILE,README.md)
	$(call CHECK_MISSING_FILE,LICENSE)
	$(call CHECK_MISSING_FILE,CODEOWNERS)
	$(call CHECK_MISSING_FILE,.github/dependabot.yml)
	$(call CHECK_MISSING_FILE,.github/workflows/ci.yml)
	$(call CHECK_MISSING_FILE,.github/workflows/dependabot_auto_merge.yml)
	$(call CHECK_MISSING_FILE,go.mod)
	$(call CHECK_MISSING_FILE,go.sum)
	$(call CHECK_MISSING_FILE,.golangci.yml)
	$(call CHECK_MISSING_FILE,Makefile)
	$(call CHECK_MISSING_FILE,Makefile-common.mk)

# Don't use upper case letter in file and directory name.
# The convention for separator in name is:
# - file: "_"
# - directory in "/cmd": "-"
# - other directory: shouldn't be separated
	@! find . -name "*.go" | (grep "[[:upper:]]" && echo "Incorrect file name case")

.PHONY: mod-update
mod-update:
	$(GO_GET) -u all
	$(MAKE) mod-tidy

.PHONY: mod-update-pierrre
mod-update-pierrre:
	GOWORK=off $(GO_LIST) -m -u -json all | jq -r 'select(.Main==null and (.Path | startswith("github.com/pierrre/")) and .Update!=null) | .Path' | xargs -I {} -t $(GO_GET) -u {}
	$(MAKE) mod-tidy

MOD_TIDY=$(GO_MOD) tidy$(VERBOSE_FLAG)
.PHONY: mod-tidy
mod-tidy:
ifeq ($(CI),true)
	$(MOD_TIDY) -diff
else
	$(MOD_TIDY)
endif

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

.PHONY: ci-tag
ci-tag:
	GOPROXY=proxy.golang.org $(GO_LIST) -x -m $(GO_MODULE)@$(GITHUB_TAG)
endif

endif # CI end
