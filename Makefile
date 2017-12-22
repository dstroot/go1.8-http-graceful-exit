#
# Variables
#

# NOTE: Simply expanded variables are defined by lines using `:='
# When a simply expanded variable is referenced, its value is substituted
# verbatim. There is another assignment operator for variables, `?='.
# This is called a conditional variable assignment operator, because it
# only has an effect if the variable is not yet defined.

# NOTE: I am not a fan of versions - people frequently forget to increment them.
# The commit ID and the buid time are more precise and automatic. However
# versions can be useful for humans so I still keep a `VERSION` file in the
# root so that anyone can clearly check the VERSION of `master`.

OWNER := dstroot
REPO := github.com
NAME := $(shell basename $(CURDIR))

PROJECT := ${REPO}/${OWNER}/${NAME}
DOCKER_NAME := ${OWNER}/${NAME}

BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_ID := $(shell git rev-parse --short HEAD 2>/dev/null || echo nosha)
VERSION := $(shell cat ./VERSION)

RELEASE_DIR := dist
GOARCH := amd64
GOOS := linux
PORT := 8000


#
# Help
#


.PHONY: help
default help:
	@echo "Usage: make <command>\n"
	@echo "The commands are:"
	@echo "   all         Alias for 'run' command."
	@echo "   gettools    Download and install Go-based build toolchain (uses go-get)."
	@echo "   clean       Clean out old builds."
	@echo "   build       Build a development version of the server. Runs dependent rules."
	@echo "   run         Run development version of the application."
	@echo "   test        Execute all development tests."
	@echo "   cover       Examine code test coverage."
	@echo "   lint        Run gometalinter against the source."
	@echo "   release     Build production release(s). Runs dependent rules."
	@echo "   docker      Build and run a local docker image."
	@echo "   dockerpush  Push the a docker image to Docker Hub."
	@echo "   minikube    Deploy the container on Kubernetes locally."
	@echo "   todo        Display all TODO's in the source."
	@echo "   docs        Display the application documentation.\n"


#
# Development
#


.PHONY: all
all: $(info Current version is $(VERSION)) run

# NOTE: Add @ to the beginning of a command to tell make not to print
# the command being executed.'

.PHONY: gettools
gettools:
	@go get -u github.com/alecthomas/gometalinter
	@go get -u github.com/golang/dep/cmd/dep
	@go get -u golang.org/x/tools/cmd/cover
	@dep ensure
	@gometalinter --install

.PHONY: clean
clean:
	@rm -f ${NAME}
	@if [ -d ${RELEASE_DIR} ]; then \
        rm -rf ${RELEASE_DIR}; \
    fi

.PHONY: build
build: clean
	@echo "Building: $(VERSION)"
	@echo "Project: $(PROJECT)"
	@echo "Name: $(NAME)"

	@go build \
		-ldflags "-s -w \
		-X ${PROJECT}/pkg/info.Version=${VERSION} \
		-X ${PROJECT}/pkg/info.Commit=${COMMIT_ID} \
		-X ${PROJECT}/pkg/info.BuildTime=${BUILD_TIME}" \
		-o ${NAME}

.PHONY: run
run: build
	@export PORT=$(PORT) && ./${NAME}


#
# Code Hygiene
#


# Go test cover with multiple packages support
.PHONY: test
test:
	@echo 'mode: atomic' > coverage.txt && go list ./... | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt' && rm coverage.tmp

# Get code coverage report
.PHONY: cover
cover: test
	@go tool cover -html=coverage.txt

# Lint all the things
.PHONY: lint
lint:
	@gometalinter --vendor ./...


#
# Release
#


# Build a Linux production executable and pack it. Since we are using a packed
# binary and Docker's scratch image we have a tiny docker container (~2mb!)
.PHONY: release
release: clean
	@echo "Releasing: $(VERSION)"
	@CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-s -w \
		-X ${PROJECT}/pkg/info.Version=${VERSION} \
		-X ${PROJECT}/pkg/info.Commit=${COMMIT_ID} \
		-X ${PROJECT}/pkg/info.BuildTime=${BUILD_TIME}" \
		-o ${RELEASE_DIR}/app

	@upx --force ${RELEASE_DIR}/*

.PHONY: push
push: release
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)


#
# Docker Section
#


# Build a docker image (assumes you have docker setup on your
# dev machine).  Use --no-cache option to force a complete rebuild.
.PHONY: dockerbuild
dockerbuild: release
	@docker build --no-cache -t $(DOCKER_NAME):latest .

# Run a docker image (assumes you have docker setup on your
# dev machine).
.PHONY: dockerrun
dockerrun: dockerbuild
	@docker stop $(DOCKER_NAME):latest || true && docker rm $(DOCKER_NAME):latest || true
	docker run -d --name ${NAME} -p ${PORT}:${PORT} \
		-e "PORT=${PORT}" \
		$(DOCKER_NAME):latest

# Push container to Docker Hub
.PHONY: dockerpush
dockerpush: dockerbuild
	@docker tag $(DOCKER_NAME):latest $(DOCKER_NAME):$(VERSION)
	@docker push $(DOCKER_NAME):latest
	@docker push $(DOCKER_NAME):$(VERSION)


#
# Minikube (Kubernetes Test)
#

minikube:
	for t in $(shell find ./kubernetes -type f -name "*.yaml"); do \
        cat $$t | \
        	gsed -E "s/\{\{(\s*)\.Release(\s*)\}\}/$(VERSION)/g" | \
        	gsed -E "s/\{\{(\s*)\.ServiceName(\s*)\}\}/$(NAME)/g"; \
        echo ---; \
    done > tmp.yaml
	kubectl apply -f tmp.yaml && rm -f tmp.yaml


#
# Misc Stuff
#


# Show any to-do items per file.
.PHONY: todo
todo:
	@grep \
	--exclude-dir=vendor \
	--exclude-dir=node_modules \
	--exclude=Makefile \
	--text \
	--color \
	-nRo -E ' TODO.*|SkipNow|nolint:.*' .

# Show documentation
.PHONY: docs
docs:
	@godoc $(shell PWD)


# NOTE: By default, Makefile targets are "file targets" - they are used to build
# files from other files. However, sometimes you want your Makefile to run
# commands that do not represent physical files in the file system. Good
# examples for this are the common targets "clean" and "all". Chances are
# this isn't the case, but you may potentially have a file named clean in
# your main directory. In such a case Make will be confused because by
# default the clean target would be associated with this file and Make
# will only run it when the file doesn't appear to be up-to-date with
# regards to its dependencies. These special targets are called phony
# and you can explicitly tell Make they're not associated with files, e.g.:
#
# .PHONY: clean
# clean:
#   rm -rf *.o
# Now make clean will run as expected even if you do have a file named clean.
#
# In terms of Make, a phony target is simply a target that is always
# out-of-date, so whenever you ask make <phony_target>, it will run,
# independent from the state of the file system. Some common make
# targets that are often phony are: all, install, clean, distclean,
# TAGS, info, check.
