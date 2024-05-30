
VERSION=0.0.1

include .bingo/Variables.mk

all: build lint test

build:
	go mod tidy
	go build ./...

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

test:
	go test ./...

release: all
	hack/tag-release.sh $(VERSION)
