
all: lint test build

VERSION=0.0.3

include .bingo/Variables.mk

VERSION_TXT=pkg/build/version.txt
SWAGGER_SPEC=swagger.json
SWAGGER_CLIENT=pkg/swagger
KORREL8RCLI=./korrel8rcli

lint: $(SWAGGER_CLIENT) $(GOLANGCI_LINT)
	go mod tidy
	$(GOLANGCI_LINT) run ./...
	@if grep -q github.com/korrel8r/korrel8r go.mod; then						\
		echo "ERROR: bad dependency: remove 'github.com/korrel8r/korrel8r' from go.mod";	\
		exit 1;	\
	fi

build: $(KORREL8RCLI)
$(KORREL8RCLI): $(VERSION_TXT) $(SWAGGER_CLIENT)
	go build -o $@  ./cmd/korrel8rcli

test:
	go test -cover -race ./...
	go tool covdata percent -i pkg/cmd/_covdata

clean:
	rm -rfv $(SWAGGER_CLIENT) $(SWAGGER_SPEC) $(KORREL8RCLI)
	git clean -dfx

ifneq ($(VERSION),$(file <$(VERSION_TXT)))
.PHONY: $(VERSION_TXT) # Force update if VERSION_TXT does not match $(VERSION)
endif

$(VERSION_TXT): $(MAKEFILE_LIST)
	echo $(VERSION) > $@

$(SWAGGER_SPEC): $(KORREL8R)	## Use bingo-installed korrel8r to generate spec.
	 $(KORREL8R) web --spec $@

$(SWAGGER_CLIENT): $(SWAGGER_SPEC) $(SWAGGER) ## Generate client packages.
	@mkdir -p $@
	cd $@ && $(SWAGGER) generate -q client -f $(abspath $(SWAGGER_SPEC))
	go mod tidy
	touch $@

pre-release: all

release: pre-release
	hack/tag-release.sh $(VERSION)
