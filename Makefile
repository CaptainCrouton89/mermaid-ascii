PREFIX ?= /usr/local/bin
pkgname ?= mermaid-ascii

targets += $(PREFIX)/$(pkgname)

ifneq (,$(wildcard /usr/share/bash-completion/completions/))
  targets += /usr/share/bash-completion/completions/$(pkgname)
endif

all: build/$(pkgname) build/completions/bash build/completions/zsh build/completions/fish ## All targets

build/$(pkgname): cmd/*.go | build/
	go build -o $@

.PHONY: install
install: $(targets)

$(PREFIX)/$(pkgname): build/$(pkgname) | $(PREFIX)
	install -m 755 $< $@
/usr/share/bash-completion/completions/$(pkgname): build/completions/bash
	install -m 755 $< $@

build/completions/%: build/$(pkgname) | build/completions/
	./$< completion $(@F) > $@

%/:
	mkdir -p $@

.PHONY: clean
clean: ## Remove generated files
	$(RM) -r build

.PHONY: uninstall
uninstall: ## Remove local installation
	$(RM) $(targets)

.PHONY: test
test: ## Run the go tests
	go test ./... -v

.PHONY: docker-build
docker-build:
	docker build -t mermaid-ascii:latest .

.PHONY: docker-test
docker-test: docker-build ## Run the go tests in docker
	docker build --target test -t mermaid-ascii:test -f Dockerfile.test .

.PHONY: docker-run
docker-run: docker-build
	docker run -i mermaid-ascii:latest

dev:
	air

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[0-9a-zA-Z._-]+:.*?## / {printf "\033[36m%s\033[0m : %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		column -s ':' -t

