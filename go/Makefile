# SPDX-License-Identifier: Apache-2.0

THIS_MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
THIS_MAKEFILE_DIR  := $(patsubst %/,%,$(dir $(THIS_MAKEFILE_PATH)))
IMAGE_NAME         := docker.io/library/golang

# Invoke this makefile one of three ways:
# With no args          : runs default target that expects the running envvironment to contain the build tools
# With arg podman       : runs Make inside podman with default target 
# With arg podman-cache : like podman, but with a second volume mount for a persistent go module cache dir of .cache

# Make using host packages - the default way of building
.PHONY: all
all: tidy compile lint format test

# Make using podman - the podman image uses all target
.PHONY: podman
podman: podman-check-image
	podman run -u root --rm -it -v $(THIS_MAKEFILE_DIR):/workdir $(IMAGE_NAME) bash -c "cd /workdir; make"

# Make using podman - as above, but with a second volume mount for persistent go module cache
.PHONY: podman-cache 
podman-cache: podman-check-image podman-check-cache
	podman run -u root --rm -it -v $(THIS_MAKEFILE_DIR)/.cache:/go/pkg/mod -v $(THIS_MAKEFILE_DIR):/workdir $(IMAGE_NAME) bash -c "cd /workdir; make" 

# Check if the required podman image has already been pulled, and pull it if not
.PHONY: podman-check-image
podman-check-image:
	[ "`podman images --format "{{.Repository}}" --filter "reference=$(IMAGE_NAME)" | wc -l`" -eq 1 ] \
	|| podman pull $(IMAGE_NAME)

# Create the persistent go mod cache dir if it does not exist 
.PHONY: podman-check-cache
podman-check-cache:
	mkdir -p $(THIS_MAKEFILE_DIR)/.cache
	
.PHONY: tidy
tidy:
	go mod tidy 

.PHONY: compile
compile:
	go build ./...

.PHONY: lint
lint:
	go vet ./...

.PHONY: format
format:
	for pkg in `go list -f '{{.Dir}}' ./...`; do gofmt -s -w $${pkg}; done

.PHONY: test
test:
	testOpt="-count=$${count:-1}"; \
	[ -n "$(run)" ] && testOpt="$$testOpt -run $(run)"; \
	go test -v $$testOpt ./...
