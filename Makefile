# SPDX-License-Identifier: Apache-2.0

# Invoke this makefile one of three ways:
# With no target           : runs default target that expects the host to contain the build tools
# With docker target       : runs Make inside docker with default target
# With docker-cache target : like docker target, but with additional volume mounts for a persistent go mod and go build caches
# With podman target       : like docker target
# With podman-cache target : like docker-cache target
#
# Passing count=N for some integer N will run the unit tests N times
#
# Passing run=T for some pattern T will run the unit tests whose names match the pattern

#### User variables, it is expected these values can be overriden on make command line

# Count of how many times to run the unit tests, default is once.
count      := 1

# Pattern for which test(s) to run, default is all of them (no pattern)
run        :=

# Package for which test(s) to run, default is all packages
pkg        := ./...

#### Variables required for Makefile

# The absolute path to this Makefile
THIS_MAKEFILE_PATH    := $(abspath $(lastword $(MAKEFILE_LIST)))

# The absolute dir containing this Makefile
THIS_MAKEFILE_DIR     := $(patsubst %/,%,$(dir $(THIS_MAKEFILE_PATH)))

# The absolute dir for caching go artifacts across docker/podman builds
GO_CACHE_ROOT        := $(THIS_MAKEFILE_DIR)/.cache

# The docker image name, in case we use docker
DOCKER_IMAGE_NAME    := golang

# Path to docker, if we have it
DOCKER_PATH           := $(shell which docker 2> /dev/null || :)

# The absolute dir for caching go package dependencies pulled in by "go mod tidy"
DOCKER_GO_CACHE_PKG   := $(GO_CACHE_ROOT)/docker/pkg

# The absolute dir for caching go build artifacts generated by "go build ./..."
DOCKER_GO_CACHE_BUILD := $(GO_CACHE_ROOT)/docker/build

# The podman image name for latest golang, in case we use podman
PODMAN_IMAGE_NAME     := docker.io/library/golang

# Path to podman, if we have it
PODMAN_PATH           := $(shell which podman 2> /dev/null || :)

# The absolute dir for caching go package dependencies pulled in by "go mod tidy"
PODMAN_GO_CACHE_PKG   := $(GO_CACHE_ROOT)/podman/pkg

# The absolute dir for caching go build artifacts generated by "go build ./..."
PODMAN_GO_CACHE_BUILD := $(GO_CACHE_ROOT)/podman/build

### Targets

# Make using host packages - the default way of building (host, docker, podman)
.PHONY: all
all: vars tidy compile lint format test spdx check-doc-go depgraph .readme.html .readme.go.html coverage

# Make using docker - the docker image uses all target (host)
.PHONY: docker
docker: docker-check-image
	docker run --rm -it -v $(THIS_MAKEFILE_DIR):/workdir $(DOCKER_IMAGE_NAME) make count=$(count) run=$(run) -C /workdir

# Make using docker - as above, but with a additional volume mounts for persistent go module and go build caches (host)
.PHONY: docker-cache
docker-cache: docker-check-image $(DOCKER_GO_CACHE_PKG) $(DOCKER_GO_CACHE_BUILD)
	docker run --rm -it -v $(DOCKER_GO_CACHE_PKG):/go/pkg -v $(DOCKER_GO_CACHE_BUILD):/root/.cache -v $(THIS_MAKEFILE_DIR):/workdir $(DOCKER_IMAGE_NAME) make uid="`id -u`" gid="`id -g`" count=$(count) run=$(run) -C /workdir

# Check if the required docker image has already been pulled, and pull it if not (host)
.PHONY: docker-check-image
docker-check-image:
	[ -n "$(DOCKER_PATH)" ] || { echo "docker is not in the path"; exit 1; }; \
	[ "`docker images --format "{{.Repository}}" --filter "reference=$(DOCKER_IMAGE_NAME)" | wc -l`" -ge 1 ] \
	|| docker pull $(DOCKER_IMAGE_NAME)

# Create directories needed for go module and go build caches (host)
$(DOCKER_GO_CACHE_PKG) $(DOCKER_GO_CACHE_BUILD):
	mkdir -p $@

# Make using podman - the podman image uses all target (host)
.PHONY: podman
podman: podman-check-image
	podman run -u root --rm -it -v $(THIS_MAKEFILE_DIR):/workdir $(PODMAN_IMAGE_NAME) make count=$(count) run=$(run) -C /workdir

# Make using podman - as above, but with a additional volume mounts for persistent go module and go build caches (host)
.PHONY: podman-cache
podman-cache: podman-check-image $(PODMAN_GO_CACHE_PKG) $(PODMAN_GO_CACHE_BUILD)
	podman run -u root --rm -it -v $(PODMAN_GO_CACHE_PKG):/go/pkg/mod -v $(PODMAN_GO_CACHE_BUILD):/root/.cache -v $(THIS_MAKEFILE_DIR):/workdir $(PODMAN_IMAGE_NAME) make count=$(count) run=$(run) -C /workdir

# Check if the required podman image has already been pulled, and pull it if not (host)
.PHONY: podman-check-image
podman-check-image:
	[ -n "$(PODMAN_PATH)" ] || { echo "podman is not in the path"; exit 1; }; \
	[ "`podman images --format "{{.Repository}}" --filter "reference=$(PODMAN_IMAGE_NAME)" | wc -l`" -ge 1 ] \
	|| podman pull $(PODMAN_IMAGE_NAME)

# Create directories needed for go module and go build caches (host)
$(PODMAN_GO_CACHE_PKG) $(PODMAN_GO_CACHE_BUILD):
	mkdir -p $@

# Download any missing go packages, update go.sum (host, docker, podman)
.PHONY: tidy
tidy:
	go mod tidy
	[ -z "$(uid)" ] || chown -R $(uid):$(gid) /go/pkg/mod # Only for docker-cache
	[ \! -d /go/pkg/mod ] || chmod -R u+w /go/pkg/mod 2> /dev/null || : # Only for docker-cache and podman-cache

# Compile go code (host, docker, podman)
.PHONY: compile
compile:
	go build ./...
	[ -z "$(uid)" ] || chown -R $(uid):$(gid) /root/.cache # Only for docker-cache
	[ \! -d /root/.cache ] || chmod -R u+w /root/.cache 2> /dev/null || : # Only for docker-cache and podman-cache

# Lint go code (host, docker, podman)
.PHONY: lint
lint:
	go vet ./...

# Format go code (host, docker, podman)
.PHONY: format
format:
	for pkg in `go list -f '{{.Dir}}' ./...`; do gofmt -s -w $${pkg}; done

# Test go code (host, docker, podman)
# May pass count=N to run tests n times
# May pass run=X to run only tests that match pattern X
# May pass pkg=X to run only tests in package X
.PHONY: test
test:
	testOpt="-count=$${count:-1}"; \
	[ -z "$(run)" ] || testOpt="$$testOpt -run $(run)"; \
	go test -coverprofile=.coverage.html -v $$testOpt $(pkg)

.PHONY: coverage
coverage:
	go tool cover -html=.coverage.html

# Check that every README and .go file contains the string SPDX-License-Identifier: Apache-2.0
.PHONY: spdx
.SILENT: spdx
spdx:
		# Recursive search for all README and .go files
		for f in $$(find $(THIS_MAKEFILE_DIR) -type f \( -iname 'README*' -o -iname '*.go' \)); do \
			[ $$(grep -c "SPDX-License-Identifier: Apache-2.0" "$$f") -gt 0 ] || { \
				echo "$$f: missing SPDX-License-Identifier: Apache-2.0"; \
				exit 1; \
			} \
		done

.PHONY: check-doc-go
check-doc-go:
	for srcDir in $$(find . -type f -name '*.go' | sed -r 's,[.]/(.*)/[^/]*,\1,' | sort -u); do \
		if [ \! -f "$$srcDir/doc.go" ]; then { echo "Missing $$srcDir/doc.go"; exit 1; }; fi; \
	done

.PHONY: have-dot
have-dot:
	@which dot 2>&1 > /dev/null || echo "The Graphviz package must be installed to generate a dependency graph"

.PHONY: depgraph
depgraph: have-dot
	{ \
		mod="$$(grep module go.mod | awk '{print $$2}')"; \
		echo 'digraph dependencies {'; \
		echo 'node [shape=box]'; \
		echo '"encoding/json" [style=filled fillcolor="#87CEFA"]'; \
		echo '"encoding/json/parse" [style=filled fillcolor="#87CEFA"]'; \
		echo '"encoding/json/write" [style=filled fillcolor="#87CEFA"]'; \
		echo '"event" [style=filled fillcolor="#87CEFA"]'; \
		echo '"rest" [style=filled fillcolor="#87CEFA"]'; \
		echo '"stream" [style=filled fillcolor="#87CEFA"]'; \
		echo '"iter" [style=filled fillcolor="#E6E6FA"]'; \
		for srcDir in $$(find . -type f -name '*.go' | sed -r 's,[.]/(.*)/[^/]*,\1,' | sort -u); do \
			(cd "$$srcDir"; go list -f '{{.Imports}}') | tr -d '[]' | tr ' ' '\n' | sort | grep "$$mod" | sed -r "s,$$mod/(.*),\"$$srcDir\" -> \"\\1\","; \
		done; \
		echo "}"; \
	} | cat > .depgraph.dot
	dot -Tsvg .depgraph.dot > depgraph.svg

.PHONY: have-asciidoc
have-asciidoc:
	@which asciidoc 2>&1 > /dev/null || echo "The asciidoc package must be installed to generate the readme html file"

.readme.html: README.adoc | have-asciidoc
	asciidoc -b html -o $@ $<

.readme.go.html: README.go.adoc | have-asciidoc
	asciidoc -b html -o $@ $<

.PHONY: push
push:
	git add -A
	git commit -m Changes
	git push

# Display all vars (host, docker, podman)
.PHONY: vars
.SILENT: vars
vars:
	printf "uid =\n    $(uid)\n"
	printf "gid =\n    $(gid)\n"
	printf "count =\n    $(count)\n"
	printf "run =\n    $(run)\n"
	printf "THIS_MAKEFILE_PATH =\n    $(THIS_MAKEFILE_PATH)\n"
	printf "THIS_MAKEFILE_DIR =\n    $(THIS_MAKEFILE_DIR)\n"
	printf "GO_CACHE_ROOT =\n    $(GO_CACHE_ROOT)\n"
	printf "DOCKER_IMAGE_NAME =\n    $(DOCKER_IMAGE_NAME)\n"
	printf "DOCKER_PATH =\n    $(DOCKER_PATH)\n"
	printf "DOCKER_GO_CACHE_PKG =\n    $(DOCKER_GO_CACHE_PKG)\n"
	printf "DOCKER_GO_CACHE_BUILD \\n    $(DOCKER_GO_CACHE_BUILD)\n"
	printf "PODMAN_IMAGE_NAME =\n    $(PODMAN_IMAGE_NAME)\n"
	printf "PODMAN_PATH =\n    $(PODMAN_PATH)\n"
	printf "PODMAN_GO_CACHE_PKG =\n    $(PODMAN_GO_CACHE_PKG)\n"
	printf "PODMAN_GO_CACHE_BUILD \\n    $(PODMAN_GO_CACHE_BUILD)\n"

# Clean artifacts (host)
.PHONY: clean
clean:
	rm -rf "$(GO_CACHE_ROOT)"
