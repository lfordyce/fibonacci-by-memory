VERSION             ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
                     			cat $(CURDIR)/.version 2> /dev/null || echo v0)
COMMIT_SHA          = $(shell git rev-parse --short HEAD)
MODULE              = $(shell env GO111MODULE=on $(GO) list -m)
DATE                ?= $(shell date +%FT%T%z)
PKGS                = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS            = $(shell env GO111MODULE=on $(GO) list -f \
                      '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
                      $(PKGS))
V                   = 0
Q                   = $(if $(filter 1,$V),,@)
M                   = $(shell printf "\033[34;1m▶\033[0m")
BIN                 = $(CURDIR)/bin
GO                  = go
GO111MODULES        ?= on
GOOS                ?= darwin
GOARCH              ?= amd64
BUILDARG            ?= build
CGO_ENABLED         ?= 0
BUILD_FLAGS         =-ldflags '-X main.Version=${VERSION} -X main.BuildDate=${DATE}'
GOBUILD             = CGO_ENABLED=$(CGO_ENABLED) GOARCH=$(GOARCH) GOOS=$(GOOS) $(GO) $(BUILDARG)

PROJECT_NAME        = $(shell basename "$(PWD)")

.PHONY: help
help: ## - Show help message
	@printf "\033[32m\xE2\x9c\x93 usage: make [target]\n\n\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: fmt | $(BIN) ; $(info $(M) building executable...) @ ## Build program binary
	$Q $(GOBUILD) \
		-tags release \
		$(BUILD_FLAGS) \
		-o $(BIN)/$(basename $(MODULE)) main.go

# Tools
$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) building $(PACKAGE)...)
	$Q tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) $(GO) get $(PACKAGE) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

GOLINT = $(BIN)/golint
$(BIN)/golint: PACKAGE=golang.org/x/lint/golint

# Tests
TEST_TARGETS := test-default test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) check test tests
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
check test tests: fmt ; $(info $(M) running $(NAME:%=% )tests...) @ ## Run tests
	$Q $(GO) test $(ARGS) $(TESTPKGS)

.PHONY: lint
lint: | $(GOLINT) ; $(info $(M) running golint...) @ ## Run golint
	$Q $(GOLINT) -set_exit_status $(PKGS)

.PHONY: clean
clean: ; $(info $(M) cleaning...)	@ ## Cleanup everything
	@rm -rf $(BIN)

.PHONY: fmt
fmt: ; $(info $(M) running gofmt...) @ ## Run gofmt on all source files
	$Q $(GO) fmt $(PKGS)

.PHONY: get
get: ; $(info $(M) running go get...) @ ## Run go get for dependencies
	$Q $(GO) get -d -v -t ./...

.PHONY: docker-pull
docker-pull: ; $(info $(M) pulling latest docker images) @ ## Pull latest Docker images in preparation for build
	$Q docker pull golang:1.16-alpine

.PHONY: docker
docker:docker-pull ; $(info $(M) building docker image) @ ## Build docker image
	$(eval BUILDER_IMAGE=$(shell docker inspect --format='{{index .RepoDigests 0}}' golang:1.16-alpine))
	@docker build -f Dockerfile --build-arg "BUILDER_IMAGE=$(BUILDER_IMAGE)" -t fib_memo_api .

.PHONY: compose
compose:docker ; $(info $(M) building images and running docker-compose up...) @ ## Run docker-compose
	$Q docker-compose -f docker-compose.yaml up