
all: generate lint test

VERSION=0.0.1-dev

include .bingo/Variables.mk

KORREL8RCLI=cmd/korrel8rcli/korrel8rcli
SWAGGER_CLIENT=pkg/swagger
SWAGGER_SPEC=korrel8r-swagger.json
GOCOVERDIR=cmd/korrel8rcli/covdata

clean:
	rm -rf $(GOCOVERDIR) $(KORREL8RCLI) $(SWAGGER_CLIENT) $(SWAGGER_SPEC)

generate: $(SWAGGER_CLIENT)

build: $(KORREL8RCLI)

$(KORREL8RCLI): $(shell find -name *.go) generate
	@mkdir -p $(dir $@)
	go build -cover -o $@ ./cmd/korrel8rcli

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

test: $(KORREL8RCLI) $(KORREL8R) $(TEST_ENV)
	KORREL8R=$(KORREL8R) go test -cover -race ./...
	@echo -e "\\nAccumulated coverage from main_test"
	go tool covdata percent -i $(GOCOVERDIR)

release: all
	hack/tag-release.sh $(VERSION)

$(SWAGGER_SPEC): $(KORREL8R)
	 $(KORREL8R) web --spec $@

$(SWAGGER_CLIENT): $(SWAGGER_SPEC) $(SWAGGER) ## Generate client packages.
	@mkdir -p $@
	cd $@ && $(SWAGGER) generate -q client -f $(abspath $(SWAGGER_SPEC)) && go mod tidy
	touch $@
