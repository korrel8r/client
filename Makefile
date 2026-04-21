# NOTE: Tool dependencies (korrel8r, golangci-lint) are managed via `go tool`.
# To update a tool version:
#   go get -tool github.com/korrel8r/korrel8r/cmd/korrel8r@VERSION
# E.g.
#   go get -tool github.com/korrel8r/korrel8r/cmd/korrel8r@latest
#   go get -tool github.com/korrel8r/korrel8r/cmd/korrel8r@v0.7.6
#   go get -tool github.com/korrel8r/korrel8r/cmd/korrel8r@main

all: lint test build

VERSION=0.0.6

VERSION_TXT=pkg/build/version.txt
OPENAPI_SPEC=korrel8r-openapi.yaml
GENERATED_CLIENT=pkg/api/generated.go

GOLANGCI_LINT=go tool golangci-lint
OAPI_CODEGEN=go tool oapi-codegen

lint: generate
	go mod tidy
	$(GOLANGCI_LINT) run ./...

generate: $(VERSION_TXT) $(GENERATED_CLIENT)

build: generate
	go build  ./cmd/korrel8rcli

install: generate
	go install  ./cmd/korrel8rcli

test: $(GENERATED_CLIENT)
	go test -cover -race ./...
	go tool covdata percent -i pkg/cmd/_covdata

clean:
	rm -rfv $(GENERATED_CLIENT) korrel8rcli
	git clean -dfx

ifneq ($(VERSION),$(file <$(VERSION_TXT)))
.PHONY: $(VERSION_TXT) # Force update if VERSION_TXT does not match $(VERSION)
endif

$(VERSION_TXT): $(MAKEFILE_LIST)
	echo $(VERSION) > $@

$(GENERATED_CLIENT): $(OPENAPI_SPEC)  ## Generate client packages.
	@mkdir -p $(dir $@)
	$(OAPI_CODEGEN) -generate types,client -package api -o $@ $<
	go mod tidy

pre-release: all

release: pre-release
	hack/tag-release.sh $(VERSION)
