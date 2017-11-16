#
# Variables
#

VERSION=$(shell cat ./VERSION)
GIT_NAME=dstroot/simple-go-webserver

# As a call to `make` without any arguments leads to the execution
# of the first target found I really prefer to make sure that this
# first one is a non-destructive one that does the most simple
# desired installation. It's very common to people set it as `all`
all: $(info $(GIT_NAME) current version is => $(VERSION))

# Install all the build and lint dependencies
setup:
	@go get -u github.com/alecthomas/gometalinter
	@go get -u github.com/golang/dep/cmd/dep
	@go get -u golang.org/x/tools/cmd/cover
	@dep ensure
	@gometalinter --install

# NOTE: Add @ to the beginning of command to tell make not to print
# the command being executed.
run:
	@echo "Running: $(VERSION)"
	@go run $(shell ls -1 *.go | grep -v _test.go)

# Install just performs a normal `go install` which builds the source
# files from the package at `./` (I like to keep a `main.go` in the root
# that imports other subpackages). After installing dependencies a
# `go install` will typically always work - except if there's an OS
# limitation in the build flags (e.g, a linux-only project).
install:
	@echo "Installing: $(VERSION)"
	go install -v

# Go test cover with multiple packages support
test:
	@echo "Testing: $(VERSION)"
	@echo 'mode: atomic' > coverage.txt && go list ./... | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt' && rm coverage.tmp

# Get code coverage report
cover: test
	@go tool cover -html=coverage.txt

# Run all the linters
lint:
	@gometalinter --vendor ./...

# Create and run a docker image (assumes you have docker setup on your
# dev machine).  Use --no-cache option to force complete rebuild.
docker:
	docker build --no-cache -t $(GIT_NAME):latest . && docker run -d -p 80:8000 --name simple $(GIT_NAME):latest

# This will create a tagged git release as well as a corresponding
# tagged docker image. I really like publishing a Docker image together
# with the GitHub release because Docker makes it very simple to someone
# run your binary without having to worry about the retrieval of the binary
# and execution of it - docker already provides the necessary boundaries.
release:
	# git
	git tag -a v$(VERSION) -m "$(GIT_NAME)-v$(VERSION)"
	git push origin v$(VERSION)
	# docker
	docker build --no-cache -t $(GIT_NAME):latest .
	docker tag $(GIT_NAME):latest $(GIT_NAME):$(VERSION)
	docker push $(GIT_NAME):latest
	docker push $(GIT_NAME):$(VERSION)

# Show any to-do items per file.
todo:
	@grep \
	--exclude-dir=vendor \
	--exclude-dir=node_modules \
	--exclude-dir=public \
	--exclude=Makefile \
	--text \
	--color \
	-nRo -E ' TODO.*|SkipNow|nolint:.*' .

# Show documentation
docs:
	@godoc $(shell PWD)

# By default, Makefile targets are "file targets" - they are used to build
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
.PHONY: install test cover todo fmt release run
