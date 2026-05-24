# NOTE: Tool dependencies (korrel8r, golangci-lint) are managed via `go tool`.
# To update a tool version:
#   go get -tool github.com/korrel8r/korrel8r/cmd/korrel8r@VERSION
# E.g.
#   go get -tool github.com/korrel8r/korrel8r/cmd/korrel8r@latest
#   go get -tool github.com/korrel8r/korrel8r/cmd/korrel8r@v0.7.6
#   go get -tool github.com/korrel8r/korrel8r/cmd/korrel8r@main

all: lint doc test

VERSION=0.1.1-dev

VERSION_TXT=pkg/build/version.txt
OPENAPI_SPEC=korrel8r-openapi.json
GEN_CLIENT=pkg/api/generated.go

GENERATED=$(VERSION_TXT) $(GEN_CLIENT)

lint: $(GENERATED)
	go mod tidy
	go tool golangci-lint run ./...

generate: $(GENERATED)

build: $(GENERATED)
	go build  ./cmd/korrel8rcli

install: $(GENERATED)
	go install  ./cmd/korrel8rcli

test: $(GENERATED)
	go test -cover -race ./...
	go tool covdata percent -i pkg/cmd/_covdata

clean: ## Remove generated files, including checked-in files.
	rm -rfv $(GENERATED) $(OPENAPI_SPEC) korrel8rcli doc/public doc/content/cmd
	git clean -dfx

ifneq ($(VERSION),$(file <$(VERSION_TXT)))
.PHONY: $(VERSION_TXT) # Force update if VERSION_TXT does not match $(VERSION)
endif

$(VERSION_TXT): $(MAKEFILE_LIST)
	echo $(VERSION) > $@

$(OPENAPI_SPEC): $(go tool -n korrel8r)
	go tool korrel8r web --spec $@

$(GEN_CLIENT): $(OPENAPI_SPEC)  ## Generate client packages.
	@mkdir -p $(dir $@)
	go tool oapi-codegen -generate types,client -package api -o $@ $<
	go mod tidy

doc/public: doc/content
	go tool hugo --source doc --quiet
	@touch $@

doc/content: doc/content/cmd/_index.md

doc/content/cmd/_index.md: $(GENERATED) $(shell find cmd pkg)
	@mkdir -p $(dir $@)
	go run ./cmd/korrel8rcli doc markdown $(dir $@)
	@printf -- '---\ntitle: Commands\n---\n' > $@
	@cat $(dir $@)korrel8rcli.md >> $@ && rm  $(dir $@)korrel8rcli.md
	for f in $(dir $@)korrel8rcli_*.md; do \
		TITLE=$$(head -1 "$$f" | sed 's/^## *//'); \
		sed -i '/### SEE ALSO/,$$d' "$$f"; \
		printf -- '---\ntitle: "%s"\n---\n' "$$TITLE" | cat - "$$f" > "$$f.tmp" && mv "$$f.tmp" "$$f"; \
	done

preview: doc
	go tool hugo server --source doc --baseURL http://localhost:1313 --bind 0.0.0.0

pre-release: all

release: pre-release
	hack/tag-release.sh $(VERSION)
