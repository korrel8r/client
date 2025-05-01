# Auto generated binary variables helper managed by https://github.com/bwplotka/bingo v0.9. DO NOT EDIT.
# All tools are designed to be build inside $GOBIN.
BINGO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO     ?= $(shell which go)

# Below generated variables ensure that every time a tool under each variable is invoked, the correct version
# will be used; reinstalling only if needed.
# For example for bingo variable:
#
# In your main Makefile (for non array binaries):
#
#include .bingo/Variables.mk # Assuming -dir was set to .bingo .
#
#command: $(BINGO)
#	@echo "Running bingo"
#	@$(BINGO) <flags/args..>
#
BINGO := $(GOBIN)/bingo-v0.9.0
$(BINGO): $(BINGO_DIR)/bingo.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/bingo-v0.9.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=bingo.mod -o=$(GOBIN)/bingo-v0.9.0 "github.com/bwplotka/bingo"

GOLANGCI_LINT := $(GOBIN)/golangci-lint-v2.1.5
$(GOLANGCI_LINT): $(BINGO_DIR)/golangci-lint.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/golangci-lint-v2.1.5"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=golangci-lint.mod -o=$(GOBIN)/golangci-lint-v2.1.5 "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"

KORREL8R := $(GOBIN)/korrel8r-v0.8.0
$(KORREL8R): $(BINGO_DIR)/korrel8r.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/korrel8r-v0.8.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=korrel8r.mod -o=$(GOBIN)/korrel8r-v0.8.0 "github.com/korrel8r/korrel8r/cmd/korrel8r"

SWAGGER := $(GOBIN)/swagger-v0.31.0
$(SWAGGER): $(BINGO_DIR)/swagger.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/swagger-v0.31.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=swagger.mod -o=$(GOBIN)/swagger-v0.31.0 "github.com/go-swagger/go-swagger/cmd/swagger"

