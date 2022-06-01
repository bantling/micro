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

# Whether to use docker or podman, in case both are installed (1 = use docker, 0 = use podman)
use_docker := 1

#### Variables required for Makefile

# The podman image name for latest golang, in case we use podman
PODMAN_IMAGE_NAME    := docker.io/library/golang

# The docker image name, in case we use docker
DOCKER_IMAGE_NAME    := $(lastword $(subst /, ,$(PODMAN_IMAGE_NAME)))

# The absolute path to this Makefile
THIS_MAKEFILE_PATH   := $(abspath $(lastword $(MAKEFILE_LIST)))

# The absolute dir containing this Makefile
THIS_MAKEFILE_DIR    := $(patsubst %/,%,$(dir $(THIS_MAKEFILE_PATH)))

# Path to docker, if we have it
DOCKER_PATH          := $(shell which docker 2> /dev/null || :)

# Path to podman, if we have it
PODMAN_PATH          := $(shell which podman 2> /dev/null || :)

ifeq ($(use_docker), 1)
CONTAINER_IMAGE_NAME := $(DOCKER_IMAGE_NAME)
CONTAINER_CMD        := $(DOCKER_PATH)
CONTAINER_RUN_CMD    := $(CONTAINER_CMD) run
else
CONTAINER_IMAGE_NAME := $(PODMAN_IMAGE_NAME)
CONTAINER_CMD        := $(PODMAN_PATH)
CONTAINER_RUN_CMD    := $(CONTAINER_CMD) run -u root
endif

# Make using host packages - the default way of building
.PHONY: all
all: tidy compile lint format test

# Make using docker - the docker image uses all target
.PHONY: docker
docker: docker-check-image
	$(CONTAINER_RUN_CMD) --rm -it -v $(THIS_MAKEFILE_DIR):/workdir $(CONTAINER_IMAGE_NAME) bash -c "cd /workdir; make count=$(count) run=$(run)"

# Make using docker - as above, but with a second volume mount for persistent go module cache
.PHONY: docker-cache
docker-cache: docker-check-image docker-check-cache
	$(CONTAINER_RUN_CMD) --rm -it -v $(THIS_MAKEFILE_DIR)/.cache:/go/pkg/mod -v $(THIS_MAKEFILE_DIR):/workdir $(CONTAINER_IMAGE_NAME) bash -c "cd /workdir; make count=$(count) run=$(run)"

# Check if the required docker image has already been pulled, and pull it if not
.PHONY: docker-check-image
docker-check-image:
	[ -n "$(CONTAINER_CMD)" ] || { echo "Neither docker nor podman are in the path"; exit 1; }; \
	[ "`$(CONTAINER_CMD) images --format "{{.Repository}}" --filter "reference=$(CONTAINER_IMAGE_NAME)" | wc -l`" -eq 1 ] \
	|| $(CONTAINER_CMD) pull $(CONTAINER_IMAGE_NAME)

# Create the persistent go mod cache dir if it does not exist 
.PHONY: docker-check-cache
docker-check-cache:
	mkdir -p $(THIS_MAKEFILE_DIR)/.cache

# Download any missing go packages, update go.sum
.PHONY: tidy
tidy:
	go mod tidy 

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
	printf "PODMAN_IMAGE_NAME =\n    $(PODMAN_IMAGE_NAME)\n"
	printf "DOCKER_IMAGE_NAME =\n    $(DOCKER_IMAGE_NAME)\n"
	printf "THIS_MAKEFILE_PATH =\n    $(THIS_MAKEFILE_PATH)\n"
	printf "THIS_MAKEFILE_DIR =\n    $(THIS_MAKEFILE_DIR)\n"
	printf "DOCKER_PATH =\n    $(DOCKER_PATH)\n"
	printf "PODMAN_PATH =\n    $(PODMAN_PATH)\n"
	printf "CONTAINER_IMAGE_NAME =\n    $(CONTAINER_IMAGE_NAME)\n"
	printf "CONTAINER_CMD =\n    $(CONTAINER_CMD)\n"
	printf "CONTAINER_RUN_CMD =\n    $(CONTAINER_RUN_CMD)\n"
