
VERSION=0.0.1-dev

include .bingo/Variables.mk

export KORREL8R
export KORREL8RCLI=$(abspath bin/korrel8rcli)
export GOCOVERDIR=$(abspath covdata)

all: lint test

clean:
	rm -rf $(GOCOVERDIR) $(KORREL8RCLI) $(KORREL8R)

build: $($KORREL8RCLI)

$(KORREL8RCLI): $(shell find -name *.go)
	@mkdir -p $(dir $@) $(GOCOVERDIR)
	go build -cover -o $@ ./cmd/korrel8rcli

lint: $(GOLANGCI_LINT)
	go mod tidy
	$(GOLANGCI_LINT) run ./...

test: $(KORREL8RCLI) $(KORREL8R)
	go test -cover -race ./...
	@echo -e "\\nAccumulated coverage from main_test"
	go tool covdata percent -i _covdata

release: all
	hack/tag-release.sh $(VERSION)


$(SWAGGER_CLIENT): $(SWAGGER_SPEC) $(SWAGGER) ## Generate client packages.
	mkdir -p $@
	cd $@ && $(SWAGGER) generate -q client -f $(abspath $(SWAGGER_SPEC)) && go mod tidy
	touch $@
