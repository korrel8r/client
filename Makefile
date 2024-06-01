
all: build lint test

VERSION=0.0.2

KORREL8RCLI=cmd/korrel8rcli/korrel8rcli
VERSION_TXT=pkg/build/version.txt
SWAGGER_CLIENT=pkg/swagger
SWAGGER_SPEC=korrel8r-swagger.json
GOCOVERDIR=cmd/korrel8rcli/covdata

include .bingo/Variables.mk

build: $(KORREL8RCLI)

lint: $(GOLANGCI_LINT)
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
$(VERSION_TXT):
	echo $(VERSION) > $@

$(KORREL8RCLI): $(VERSION_TXT) $(SWAGGER_CLIENT) $(shell find -name *.go)
	@mkdir -p $(dir $@)
	@go mod tidy
	go build -cover -o $@ ./cmd/korrel8rcli

$(SWAGGER_SPEC): $(KORREL8R)
	 $(KORREL8R) web --spec $@

$(SWAGGER_CLIENT): $(SWAGGER_SPEC) $(SWAGGER) ## Generate client packages.
	@mkdir -p $@
	cd $@ && $(SWAGGER) generate -q client -f $(abspath $(SWAGGER_SPEC)) && go mod tidy
	touch $@

release: all
	hack/tag-release.sh $(VERSION)
