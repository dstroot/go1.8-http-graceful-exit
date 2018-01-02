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
# versions can be useful for humans, so I still keep a `VERSION` file in the
# root so that anyone can clearly check the VERSION of `master`.

OWNER := dstroot
REPO := github.com
NAME := $(shell basename $(CURDIR))

PROJECT := ${REPO}/${OWNER}/${NAME}
DOCKER_NAME := ${OWNER}/${NAME}

BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_ID := $(shell git rev-parse --short HEAD 2>/dev/null || echo nosha)
VERSION := $(shell cat ./VERSION)

GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
	COMMIT_ID := $(COMMIT_ID)-dirty
endif

RELEASE_DIR := dist
GOARCH := amd64
GOOS := linux
PORT := 8000


#
# Help
#


.PHONY: help
help: ## Display help.
	@echo "Usage: make <command> \n"
	@echo "Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-30s\033[0m %s\n", $$1, $$2}'


#
# Development
#


.PHONY: all
all: $(info Current version is $(VERSION)-$(COMMIT_ID)) clean build run ## Runs "clean", "build", and "run".

# NOTE: Add @ to the beginning of a command to tell make not to print
# the command being executed.'

.PHONY: gettools
gettools: ## Download and install Go-based build toolchain (uses go-get).
	@go get -u github.com/alecthomas/gometalinter
	@go get -u github.com/golang/dep/cmd/dep
	@go get -u golang.org/x/tools/cmd/cover
	@dep ensure
	@gometalinter --install

.PHONY: clean
clean: ## Cleanup any build binaries or packages.
	@echo "+ $@"
	@$(RM) $(NAME)
	@if [ -d ${RELEASE_DIR} ]; then \
		$(RM) -r ${RELEASE_DIR}; \
	fi

.PHONY: build
build: $(NAME) ## Builds a dynamic executable or package.

$(NAME): *.go VERSION
	@echo "+ build $@ v$(VERSION)"
	@go build \
		-ldflags "-s -w \
		-X ${PROJECT}/pkg/info.Version=${VERSION} \
		-X ${PROJECT}/pkg/info.Commit=${COMMIT_ID} \
		-X ${PROJECT}/pkg/info.BuildTime=${BUILD_TIME}" \
		-o $(NAME) .

.PHONY: run
run: build ## Run the project executable.
	@echo "+ $@"
	@export PORT=$(PORT) && ./${NAME}


#
# Code Hygiene
#


.PHONY: test
test: ## Run the go tests, including /pkg tests.
	@echo "+ $@"
	@echo 'mode: atomic' > coverage.txt && go list ./... | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt' && rm coverage.tmp

.PHONY: cover
cover: test ## Get code coverage report.
	@echo "+ $@"
	@go tool cover -html=coverage.txt

.PHONY: lint
lint: ## Lint all the things.
	@echo "+ $@"
	@gometalinter --vendor ./...


#
# Release
#


.PHONY: release
release: $(RELEASE_DIR) ## Build a Linux production executable and pack it.

$(RELEASE_DIR): *.go VERSION
	@echo "+ release $@"

	@echo "Releasing: $(VERSION)"
	@echo "Project: $(PROJECT)"
	@echo "Name: $(NAME)"

	@CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-a -tags "static_build netgo" \
		-ldflags "-s -w \
		-X ${PROJECT}/pkg/info.Version=${VERSION} \
		-X ${PROJECT}/pkg/info.Commit=${COMMIT_ID} \
		-X ${PROJECT}/pkg/info.BuildTime=${BUILD_TIME} \
		-extldflags -static" \
		-o ${RELEASE_DIR}/app .

	@upx --force ${RELEASE_DIR}/*

	openssl md5 $(RELEASE_DIR)/app > $(RELEASE_DIR)/$(NAME)-$(VERSION).md5;
	openssl sha256 $(RELEASE_DIR)/app > $(RELEASE_DIR)/$(NAME)-$(VERSION).sha256;

.PHONY: push
push: dockerpush ## Push new version tag to Github and push tagged Docker image.
	@echo "+ $@"
	git tag -sa $(VERSION) -m "$(VERSION)"
	git push origin $(VERSION)


#
# Docker Section
#


# Note: Use --no-cache option to force a complete rebuild.
.PHONY: dockerbuild
dockerbuild: release ## Build a Docker image (assumes you have docker installed).
	@echo "+ $@"
	@docker build --no-cache -t $(DOCKER_NAME):latest .

.PHONY: dockerrun
dockerrun: dockerbuild ## Run docker image locally (assumes you have docker installed).
	@echo "+ $@"
	@docker stop $(DOCKER_NAME):latest || true && docker rm $(DOCKER_NAME):latest || true
	docker run -d --name ${NAME} -p ${PORT}:${PORT} \
		-e "PORT=${PORT}" \
		$(DOCKER_NAME):latest

.PHONY: dockerpush
dockerpush: dockerbuild ## Deploy the project to Docker Hub.
	@echo "+ $@"
	@docker tag $(DOCKER_NAME):latest $(DOCKER_NAME):$(VERSION)
	@docker push $(DOCKER_NAME):latest
	@docker push $(DOCKER_NAME):$(VERSION)


#
# Minikube (Kubernetes Test)
#

.PHONY: minikube
minikube: ## Deploy the project to minikube.
	@echo "+ $@"
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


.PHONY: todo
todo: ## Display any "TODOs" in the source code.
	@echo "+ $@"
	@grep \
	--exclude-dir=public \
	--exclude-dir=vendor \
	--exclude-dir=node_modules \
	--exclude=Makefile \
	--text \
	--color \
	-nRo -E ' TODO.*|SkipNow|nolint:.*' .

.PHONY: docs
docs: ## Display project docs.
	@echo "+ $@"
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
#
# https://github.com/jessfraz/weather/blob/master/Makefile
