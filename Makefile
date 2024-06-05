
all: build lint test

VERSION=0.0.3

KORREL8RCLI=cmd/korrel8rcli/korrel8rcli
VERSION_TXT=pkg/build/version.txt
SWAGGER_SPEC=swagger.json
SWAGGER_CLIENT=pkg/swagger
GOCOVERDIR=cmd/korrel8rcli/covdata

include .bingo/Variables.mk

build: $(KORREL8RCLI)

lint: $(SWAGGER_CLIENT) $(GOLANGCI_LINT)
	go mod tidy
	@if grep -q github.com/korrel8r/korrel8r go.mod; then echo "ERROR: go.mod contains 'github.com/korrel8r/korrel8r'; remove this dependency."; exit 1; fi
	$(GOLANGCI_LINT) run ./...

test: $(KORREL8RCLI) $(KORREL8R)
	KORREL8R=$(KORREL8R) go test -cover -race ./...
	@echo -e "\\nAccumulated coverage from main_test"
	go tool covdata percent -i $(GOCOVERDIR)

clean:
	rm -rf $(GOCOVERDIR) $(KORREL8RCLI) $(SWAGGER_CLIENT) $(SWAGGER_SPEC)

ifneq ($(VERSION),$(file <$(VERSION_TXT)))
.PHONY: $(VERSION_TXT) # Force update if VERSION_TXT does not match $(VERSION)
endif

$(KORREL8RCLI): $(VERSION_TXT) $(SWAGGER_CLIENT) $(shell find -name *.go)
	@mkdir -p $(dir $@)
	go build -cover -o $@ ./cmd/korrel8rcli

$(VERSION_TXT): $(MAKEFILE_LIST)
	echo $(VERSION) > $@

$(SWAGGER_SPEC): $(KORREL8R)
	 $(KORREL8R) web --spec $@

$(SWAGGER_CLIENT): $(SWAGGER_SPEC) $(SWAGGER) ## Generate client packages.
	@mkdir -p $@
	cd $@ && $(SWAGGER) generate -q client -f $(abspath $(SWAGGER_SPEC))
	go mod tidy
	touch $@

pre-release: all

release: pre-release
	hack/tag-release.sh $(VERSION)
