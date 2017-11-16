# # Retrieve the `golang:alpine` image to provide us the
# # necessary Golang tooling for building Go binaries.
# # Here I retrieve the `alpine`-based just for the
# # convenience of using a smaller image.
# FROM golang:alpine as builder
#
# # Add the `main` file that is really the only golang
# # file under the root directory that matters for the
# # build
# ADD ./main.go /go/src/github.com/cirocosta/l7/main.go
#
# # Add all the files from the packages that I own
# ADD ./lib /go/src/github.com/cirocosta/l7/lib
#
# # Add vendor dependencies (committed or not)
# # I typically commit the vendor dependencies as it
# # makes the final build more reproducible and less
# # dependant on dependency managers.
# ADD ./vendor /go/src/github.com/cirocosta/l7/vendor
#
# # 0.    Set some shell flags like `-e` to abort the
# #       execution in case of any failure (useful if we
# #       have many ';' commands) and also `-x` to print to
# #       stderr each command already expanded.
# # 1.    Get into the directory with the golang source code
# # 2.    Perform the go build with some flags to make our
# #       build produce a static binary (CGO_ENABLED=0 and
# #       the `netgo` tag).
# # 3.    copy the final binary to a suitable location that
# #       is easy to reference in the next stage
# RUN set -ex && \
#   cd /go/src/github.com/cirocosta/l7 && \
#   CGO_ENABLED=0 go build \
#         -tags netgo \
#         -v -a \
#         -ldflags '-extldflags "-static"' && \
#   mv ./l7 /usr/bin/l7
#
# # Create the second stage with the most basic that we need - a
# # busybox which contains some tiny utilities like `ls`, `cp`,
# # etc. When we do this we'll end up dropping any previous
# # stages (defined as `FROM <some_image> as <some_name>`)
# # allowing us to start with a fat build image and end up with
# # a very small runtime image. Another common option is using
# # `alpine` so that the end image also has a package manager.
# # You might want to end up with a FROM alpine instead of FROM
# # busybox and then run apk add --update ca-certificates if your
# # binary needs to perform requests to HTTPS endpoints - just
# # using busybox will lead to an image that doesn’t contain root
# # CA certificates which would make HTTPS requests fail.
# FROM busybox
#
# # Retrieve the binary from the previous stage
# COPY --from=builder /usr/bin/l7 /usr/local/bin/l7
#
# # Set the binary as the entrypoint of the container
# ENTRYPOINT [ "l7" ]

# # Anyone using a language capable of producing statically linked
# binaries has the opportunity to package based on “scratch” (the
# empty image). There are quite a few reasons to do this:
#   * avoid obscure licensing issues
#   * reduced attack surface
#   * explicit dependency management

# However one constant pain is dealing with CA root certificates.
# Images based on an image like alpine can use built in package
# management to fetch updated CA certificates provided by the
# distribution. There is no package management for scratch.
# You have to bring your own CA certificate bundle.

# Multi-stage builds make it easier to borrow automations and
# artifacts from other images. The example below uses the alpine
# package manager to fetch the current ca-certificates package
# and later copy the downloaded artifact into a different image
# based on scratch.
FROM golang:latest as builder
# RUN apk --update add ca-certificates

# Build our binary
RUN CGO_ENABLED=0 go get -a -ldflags '-s' github.com/dstroot/simple-go-webserver

# Build the final container. Use scratch as the smallest possible container
FROM scratch
# Add in certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# Add the binary
COPY --from=builder /go/bin/simple-go-webserver .
ADD ./public /public
ADD ./templates /templates
# Run it
EXPOSE 8000
CMD ["./simple-go-webserver"]
