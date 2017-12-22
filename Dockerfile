# Multi-stage builds make it easier to borrow automations and
# artifacts from other images. The example below uses the alpine
# package manager to fetch the current ca-certificates package
# and later copy the downloaded artifact into a different image
# based on scratch.
FROM alpine:3.7 as builder
RUN apk --update add ca-certificates

# Build the final container. Use scratch: the smallest possible container
FROM scratch

# Add in certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# set the port
ENV PORT 8000
EXPOSE $PORT

COPY dist/app /
COPY templates /templates
COPY public /public

CMD ["./app"]
