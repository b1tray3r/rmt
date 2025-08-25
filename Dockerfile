FROM golang:1.24-alpine AS builder

# Build arguments for metadata
ARG VERSION=dev
ARG BUILD_DATE
ARG VCS_REF

LABEL org.opencontainers.image.source="https://github.com/b1tray3r/rmt"
LABEL org.opencontainers.image.description="RMT - Redmine Time tracking tool"
LABEL org.opencontainers.image.url="https://github.com/b1tray3r/rmt"
LABEL org.opencontainers.image.documentation="https://github.com/b1tray3r/rmt"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.version="$VERSION"
LABEL org.opencontainers.image.revision="$VCS_REF"
LABEL org.opencontainers.image.created="$BUILD_DATE"

RUN apk update && \
    apk add --no-cache \
        gcc \
        musl-dev \
        ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Build with security flags and version information
ENV CGO_ENABLED=1
RUN go build \
    -ldflags="-s -w -extldflags '-static' -X main.version=$VERSION -X main.buildDate=$BUILD_DATE -X main.gitCommit=$VCS_REF" \
    -a -installsuffix cgo \
    -o rmt .

FROM alpine:3.22

RUN apk update && \
    apk add --no-cache ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

RUN addgroup -g 10001 -S nonroot && \
    adduser -u 10001 -S -G nonroot -D -H -s /sbin/nologin nonroot

WORKDIR /app
RUN chown nonroot:nonroot /app

COPY --from=builder --chown=nonroot:nonroot /build/rmt /usr/local/bin/rmt
COPY --chown=nonroot:nonroot ./VERSION /VERSION

# Switch to non-root user before setting up runtime
USER nonroot

ENTRYPOINT ["/usr/local/bin/rmt"]
