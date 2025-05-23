# NOTE: This Makefile installs korrel8r using bingo.
# The bingo-installed version is used for testing and to generate the REST API swagger spec.
# To update the version of korrel8r:
#   bingo get korrel8r@VERSION
# E.g.
#   bingo get korrel8r@latest # latest released version
#   bingo get korrel8r@v0.7.6 # specific version.
#   bingo get korrel8r@main   # latest development snapshot on main.
# To see what version is being used:
#   bingo list
#

all: lint test build

VERSION=0.0.6

VERSION_TXT=pkg/build/version.txt
SWAGGER_SPEC=swagger.json
SWAGGER_CLIENT=pkg/swagger

include .bingo/Variables.mk

lint: generate $(GOLANGCI_LINT)
	go mod tidy
	$(GOLANGCI_LINT) run ./...
	@if grep -q github.com/korrel8r/korrel8r go.mod; then						\
		echo "ERROR: bad dependency: remove 'github.com/korrel8r/korrel8r' from go.mod";	\
		exit 1;	\
	fi

generate: $(VERSION_TXT) $(SWAGGER_CLIENT)

build: generate
	go build  ./cmd/korrel8rcli

install: generate
	go install  ./cmd/korrel8rcli

export KORREL8R
test: $(SWAGGER_CLIENT) $(KORREL8R)
	go test -cover -race ./...
	go tool covdata percent -i pkg/cmd/_covdata

clean:
	rm -rfv $(SWAGGER_CLIENT) korrel8rcli
	git clean -dfx

ifneq ($(VERSION),$(file <$(VERSION_TXT)))
.PHONY: $(VERSION_TXT) # Force update if VERSION_TXT does not match $(VERSION)
endif

$(VERSION_TXT): $(MAKEFILE_LIST)
	echo $(VERSION) > $@

$(SWAGGER_CLIENT): $(SWAGGER_SPEC) $(SWAGGER) ## Generate client packages.
	@mkdir -p $@
	cd $@ && $(SWAGGER) generate -q client -f $(abspath $(SWAGGER_SPEC))
	go mod tidy
	touch $@

$(SWAGGER_SPEC): $(KORREL8R)
	$(KORREL8R) web --spec $@

pre-release: all

release: pre-release
	hack/tag-release.sh $(VERSION)
