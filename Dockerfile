# -----------------------------------
# 1. Build Stage: Compile the Go app
# -----------------------------------
ARG GO_VERSION=1.24.1
ARG CONTAINER_PORT=8080
ARG MOUNT_PATH=/app/keys

FROM golang:${GO_VERSION}-alpine AS builder

# Install required tools (e.g., certificates, Git)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy module files first (use caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code and compile
COPY . .
RUN go build -o oauth2-server ./cmd/server

# ----------------------------------------
# 2. Runtime Stage: Create a minimal image
# ----------------------------------------
FROM alpine:latest

ARG CONTAINER_PORT
ARG MOUNT_PATH

# Install CA certificates if TLS is used (e.g., for JWKS over HTTPS)
RUN apk add --no-cache ca-certificates

# Copy binary - no keys
COPY --from=builder /app/oauth2-server /oauth2-server

# Runtime expects /app/keys – will be mounted
VOLUME ["${MOUNT_PATH}"]

EXPOSE ${CONTAINER_PORT}

# Define entrypoint
ENTRYPOINT ["/oauth2-server"]