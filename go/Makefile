# SPDX-License-Identifier: Apache-2.0

# Invoke this makefile one of three ways:
# With no target           : runs default target that expects the host to contain the build tools
# With docker target       : runs Make inside docker or podman with default target
# With docker-cache target : like docker target, but with a second volume mount for a persistent go module cache dir of .cache
#
# If docker does not exist, podman will be used instead

# If docker and podman both exist:
# - by default, docker will be used
# - passing use_docker=0 will use podman

# If neither docker nor podman exist:
# - an error occurs if you call the docker targets
# - no error occurs if you do not call the docker targets

# Passing count=N for some integer N will run the unit tests N times

# Passing run=T for some pattern T will run the unit tests whose names match the pattern

#### User variables, it is expected these values can be overriden on make command line

# Count of how many times to run the unit tests, default is once.
count      := 1

# Pattern for which test(s) to run, default is all of them (no pattern)
run        :=

#### Variables required for Makefile

# The docker image name, in case we use docker
DOCKER_IMAGE_NAME    := golang

# The podman image name for latest golang, in case we use podman
PODMAN_IMAGE_NAME    := docker.io/library/golang

# The absolute path to this Makefile
THIS_MAKEFILE_PATH   := $(abspath $(lastword $(MAKEFILE_LIST)))

# The absolute dir containing this Makefile
THIS_MAKEFILE_DIR    := $(patsubst %/,%,$(dir $(THIS_MAKEFILE_PATH)))

# Path to docker, if we have it
DOCKER_PATH          := $(shell which docker 2> /dev/null || :)

# Path to podman, if we have it
PODMAN_PATH          := $(shell which podman 2> /dev/null || :)

# Make using host packages - the default way of building
.PHONY: all
all: tidy compile lint format test

# Make using docker - the docker image uses all target
.PHONY: docker
docker: docker-check-image
	docker run --rm -it -v $(THIS_MAKEFILE_DIR):/workdir $(DOCKER_IMAGE_NAME) make -C /workdir count=$(count) run=$(run)

# Make using docker - as above, but with a second volume mount for persistent go module cache
.PHONY: docker-cache
docker-cache: docker-check-image $(THIS_MAKEFILE_DIR)/.cache
	docker run --rm -it -v $(THIS_MAKEFILE_DIR)/.cache:/go/pkg/mod -v $(THIS_MAKEFILE_DIR):/workdir $(DOCKER_IMAGE_NAME) make -C /workdir id="`id -u`" count=$(count) run=$(run)

# Check if the required docker image has already been pulled, and pull it if not
.PHONY: docker-check-image
docker-check-image:
	[ -n "$(DOCKER_PATH)" ] || { echo "docker is not in the path"; exit 1; }; \
	[ "`docker images --format "{{.Repository}}" --filter "reference=$(DOCKER_IMAGE_NAME)" | wc -l`" -eq 1 ] \
	|| docker pull $(DOCKER_IMAGE_NAME)

# Make using podman - the podman image uses all target
.PHONY: podman
podman: podman-check-image
	podman run -u root --rm -it -v $(THIS_MAKEFILE_DIR):/workdir $(PODMAN_IMAGE_NAME) bash -c "cd /workdir; make count=$(count) run=$(run)"

# Make using podman - as above, but with a second volume mount for persistent go module cache
.PHONY: podman-cache
podman-cache: podman-check-image $(THIS_MAKEFILE_DIR)/.cache
	podman run --rm -it -v $(THIS_MAKEFILE_DIR)/.cache:/go/pkg/mod -v $(THIS_MAKEFILE_DIR):/workdir $(PODMAN_IMAGE_NAME) bash -c "echo $$HOME; cd /workdir; make count=$(count) run=$(run)"

# Check if the required podman image has already been pulled, and pull it if not
.PHONY: podman-check-image
podman-check-image:
	[ -n "$(PODMAN_PATH)" ] || { echo "podman is not in the path"; exit 1; }; \
	[ "`podman images --format "{{.Repository}}" --filter "reference=$(PODMAN_IMAGE_NAME)" | wc -l`" -eq 1 ] \
	|| podman pull $(PODMAN_IMAGE_NAME)

# Create the persistent go mod cache dir if it does not exist 
$(THIS_MAKEFILE_DIR)/.cache:
	mkdir -p $(THIS_MAKEFILE_DIR)/.cache

# Download any missing go packages, update go.sum
.PHONY: tidy
tidy:
	go mod tidy
	if [ -d $(THIS_MAKEFILE_DIR)/.cache ]; then \
		[ -z "${id}" ] || chown -R ${id} $(THIS_MAKEFILE_DIR)/.cache; \
		chmod -R u+w $(THIS_MAKEFILE_DIR)/.cache; \
	fi

# Compile go code
.PHONY: compile
compile:
	go build ./...

# Lint go code
.PHONY: lint
lint:
	go vet ./...

# Format go code
.PHONY: format
format:
	for pkg in `go list -f '{{.Dir}}' ./...`; do gofmt -s -w $${pkg}; done

# Test go code
# May pass count=N to run tests n times
.PHONY: test
test:
	testOpt="-count=$${count:-1}"; \
	[ -n "$(run)" ] && testOpt="$$testOpt -run $(run)"; \
	go test -v $$testOpt ./...

.PHONY: vars
.SILENT: vars
vars:
	printf "DOCKER_IMAGE_NAME =\n    $(DOCKER_IMAGE_NAME)\n"
	printf "PODMAN_IMAGE_NAME =\n    $(PODMAN_IMAGE_NAME)\n"
	printf "THIS_MAKEFILE_PATH =\n    $(THIS_MAKEFILE_PATH)\n"
	printf "THIS_MAKEFILE_DIR =\n    $(THIS_MAKEFILE_DIR)\n"
	printf "DOCKER_PATH =\n    $(DOCKER_PATH)\n"
	printf "PODMAN_PATH =\n    $(PODMAN_PATH)\n"
	printf "HOME =\n    $(HOME)\n"
	printf "ID =\n    `id`\n"
